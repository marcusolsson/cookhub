package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

func (s *Server) apiAddRepository(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Slug string `json:"slug"`
		Ref  string `json:"ref"`
	}

	json.NewDecoder(req.Body).Decode(&body)

	if body.Ref == "" {
		body.Ref = "HEAD"
	}

	segments := strings.Split(body.Slug, "/")

	if len(segments) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.db.AddRepository(context.Background(), db.AddRepositoryParams{
		Ref:      body.Ref,
		Url:      body.Slug,
		Provider: segments[0],
		Owner:    segments[1],
		RepoName: segments[2],
		Slug:     segments[1] + "/" + segments[2],
	}); err != nil {
		s.logger.Error("Failed to add repository", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) apiImportRepos(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	repos, err := s.db.ListRepositories(ctx)
	if err != nil {
		s.logger.Error("Failed to list repositories", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.logger.Info("Import started", "repo_count", len(repos))

	numWorkers := 3
	ch := make(chan db.Repository, len(repos))

	var wg sync.WaitGroup

	var processedRepos int
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for repo := range ch {
				ctxlog := s.logger.With("repository_id", repo.ID, "repository", repo.Url)

				sha, err := s.gh.GetLatestCommitSHA(ctx, repo.Owner, repo.RepoName, repo.Ref)
				if err != nil {
					ctxlog.Error("Failed to get latest commit SHA from repo", "error", err)
					continue
				}

				if sha == "" {
					ctxlog.Error("Empty commit SHA")
					continue
				}

				ctxlog = ctxlog.With("commit_sha", sha)

				last_sha, err := s.db.GetCommitShaByRepo(ctx, repo.ID)
				if err != nil && !errors.Is(err, pgx.ErrNoRows) {
					ctxlog.Error("Failed to get last commit SHA from database", "error", err)
					continue
				}

				if sha == last_sha {
					ctxlog.Info("Skipping up-to-date repository")
					continue
				}

				ctxlog.Info("Ingesting files")

				if err := s.withIngestionRun(ctx, repo, sha, s.ingestFiles); err != nil {
					ctxlog.Info("Failed to ingest files", "error", err.Error())
				}

				processedRepos++
			}
		}()
	}

	for _, repo := range repos {
		ch <- repo
	}
	close(ch)

	wg.Wait()

	s.logger.Info("Import finished", "imported_repos", processedRepos)
}

func (s *Server) withIngestionRun(
	ctx context.Context,
	repo db.Repository,
	commitSHA string,
	fn func(ctx context.Context, qs *db.Queries, repo db.Repository, runID, sha string) (int, error),
) error {
	runID, err := s.db.CreateIngestionRun(ctx, db.CreateIngestionRunParams{
		RepoID:    repo.ID,
		RepoRef:   repo.Ref,
		CommitSha: commitSHA,
	})
	if err != nil {
		return err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	numFiles, err := fn(ctx, s.db.WithTx(tx), repo, runID, commitSHA)
	if err != nil {
		_ = tx.Rollback(ctx)

		return s.db.MarkRunAsFailed(ctx, db.MarkRunAsFailedParams{
			ID:                  runID,
			FilesProcessedCount: pgtype.Int4{Int32: int32(numFiles), Valid: true},
			ErrorMessage:        pgtype.Text{String: err.Error(), Valid: true},
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return s.db.MarkRunAsFailed(ctx, db.MarkRunAsFailedParams{
			ID:                  runID,
			FilesProcessedCount: pgtype.Int4{Int32: int32(numFiles), Valid: true},
			ErrorMessage:        pgtype.Text{String: err.Error(), Valid: true},
		})
	}

	return s.db.MarkRunAsCompleted(ctx, db.MarkRunAsCompletedParams{
		ID:                  runID,
		FilesProcessedCount: pgtype.Int4{Int32: int32(numFiles), Valid: true},
	})
}

func (s *Server) ingestFiles(
	ctx context.Context,
	qs *db.Queries,
	repo db.Repository,
	runID, sha string,
) (int, error) {
	tempDir, err := os.MkdirTemp("", "cooklang-")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tempDir)

	zipPath, err := downloadZipBall(ctx, repo.Owner+"/"+repo.RepoName, sha, tempDir)
	if err != nil {
		return 0, err
	}

	var runParams []db.ImportFilesParams
	for file, err := range readFilesFromZip(zipPath) {
		if err != nil {
			return 0, err
		}

		hash := sha256.Sum256(file.Content)

		runParams = append(runParams, db.ImportFilesParams{
			IngestionRunID: runID,
			Path:           file.Name,
			Basename:       filepath.Base(file.Name),
			Extension:      filepath.Ext(file.Name),
			Stem:           strings.TrimSuffix(filepath.Base(file.Name), filepath.Ext(file.Name)),
			Content:        string(file.Content),
			SizeBytes:      int32(len(file.Content)),
			Hash:           hash[:],
		})
	}

	_, err = qs.ImportFiles(ctx, runParams)
	if err != nil {
		return len(runParams), err
	}

	return len(runParams), nil
}
