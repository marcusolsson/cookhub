package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

func (s *server) apiIndexRepo(w http.ResponseWriter, req *http.Request) {
	var (
		provider = chi.URLParam(req, "provider")
		owner    = chi.URLParam(req, "owner")
		name     = chi.URLParam(req, "name")
		fullName = fmt.Sprintf("%s/%s", owner, name)
	)

	ctx := req.Context()

	ctxlog := s.logger.With("provider", provider, "owner", owner, "name", name)

	_, body, err := s.ghClient.GetRepository(ctx, owner, name)
	if err != nil {
		ctxlog.Error("Failed to get metadata for repository", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sha, err := s.ghClient.GetLatestCommitSHA(ctx, owner, name, "HEAD")
	if err != nil {
		ctxlog.Error("Failed to get latest commit SHA from repo", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sha == "" {
		ctxlog.Error("Commit SHA was empty")
		return
	}

	ctxlog = ctxlog.With("branch", "HEAD", "commit_sha", sha)

	jobID, err := s.db.CreateJob(ctx, db.CreateJobParams{Slug: fullName, CommitSha: sha})
	if err != nil {
		ctxlog.Error("Failed to create job", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctxlog = ctxlog.With("job_id", jobID)

	go func() {
		ctx := context.Background()

		tx, err := s.pool.Begin(ctx)
		if err != nil {
			ctxlog.Error("Failed to begin transaction", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
				ctxlog.Error("Failed to roll back transaction", "error", err)
			}
		}()

		qtx := s.db.WithTx(tx)

		if err := s.ingestRepoMetadata(ctx, qtx, jobID, provider, owner, name, body); err != nil {
			ctxlog.Error("Failed to ingest repo metadata", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		numFiles, err := s.ingestFiles(ctx, qtx, jobID, fullName, sha)
		if err != nil {
			ctxlog.Error("Failed to ingest files", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			ctxlog.Error("Failed to commit transaction", "error", err)
			if statusErr := s.db.SetJobStatus(ctx, db.SetJobStatusParams{
				ID:     jobID,
				Status: db.StatusEnumFailed,
			}); statusErr != nil {
				ctxlog.Error("Failed to update job status", "error", statusErr)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := s.db.SetJobStatus(ctx, db.SetJobStatusParams{
			ID:     jobID,
			Status: db.StatusEnumCompleted,
		}); err != nil {
			ctxlog.Error("Failed to update job status", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctxlog.Info("Successfully indexed repository", "num_files", numFiles)
	}()

	ctxlog.Info("Started indexing repository")
}

func (s *server) ingestRepoMetadata(
	ctx context.Context,
	qs *db.Queries,
	jobID, provider, owner, name string,
	body []byte,
) error {
	return qs.AddRepositoryMetadata(ctx, db.AddRepositoryMetadataParams{
		JobID:    jobID,
		Provider: provider,
		Owner:    owner,
		Name:     name,
		Response: body,
	})
}

func (s *server) ingestFiles(
	ctx context.Context,
	qs *db.Queries,
	jobID, slug, sha string,
) (int, error) {
	tempDir, err := os.MkdirTemp("", "cooklang-")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tempDir)

	zipPath, err := downloadZipBall(ctx, slug, sha, tempDir)
	if err != nil {
		return 0, err
	}

	var params []db.CreateFileParams
	for file, err := range readFilesFromZip(zipPath) {
		if err != nil {
			return 0, err
		}

		params = append(params, db.CreateFileParams{
			JobID:   jobID,
			Name:    file.Name,
			Content: string(file.Content),
		})
	}

	_, err = qs.CreateFile(ctx, params)
	if err != nil {
		return 0, err
	}

	return len(params), nil
}

func (s *server) apiListMetadata(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	files, err := s.db.GetAllFiles(ctx)
	if err != nil {
		s.logger.Error("Failed to get files from database", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fields := map[string]int{}

	for _, file := range files {
		recipe, err := parseCooklangRecipe(file.Content)
		if err != nil {
			continue
		}

		for key := range recipe.Metadata {
			if _, ok := fields[key]; !ok {
				fields[key] = 0
			}

			fields[key]++
		}
	}

	for key, count := range fields {
		s.logger.Info(key, "count", count)
	}
}
