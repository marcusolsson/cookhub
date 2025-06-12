package main

import (
	"context"
	"crypto/sha256"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

func (s *Server) apiImportRepos(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	repos, err := s.db.ListRepositories(ctx)
	if err != nil {
		s.logger.Error("Failed to get repositories", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, repo := range repos {
		fullName := repo.Owner + "/" + repo.RepoName
		ctxlog := s.logger.With("repository_id", repo.ID)

		sha, err := s.gh.GetLatestCommitSHA(ctx, repo.Owner, repo.RepoName, "HEAD")
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

		runID, err := s.db.CreateIngestionRun(ctx, db.CreateIngestionRunParams{
			RepositoryID: repo.ID,
			Branch:       repo.Branch,
			CommitSha:    sha,
		})
		if err != nil {
			ctxlog.Error("Failed to create run", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

			numFiles, err := s.ingestFiles(ctx, qtx, runID, fullName, sha)
			if err != nil {
				_ = qtx.MarkRunAsFailed(ctx, db.MarkRunAsFailedParams{
					ID:                  runID,
					FilesProcessedCount: pgtype.Int4{Int32: int32(numFiles), Valid: true},
					ErrorMessage:        pgtype.Text{String: err.Error(), Valid: true},
				})

				ctxlog.Error("Failed to ingest files", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := tx.Commit(ctx); err != nil {
				ctxlog.Error("Failed to commit transaction", "error", err)
				if statusErr := s.db.MarkRunAsFailed(ctx, db.MarkRunAsFailedParams{
					ID:                  runID,
					FilesProcessedCount: pgtype.Int4{Int32: int32(numFiles), Valid: true},
					ErrorMessage:        pgtype.Text{String: err.Error(), Valid: true},
				}); statusErr != nil {
					ctxlog.Error("Failed to update job status", "error", statusErr)
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if statusErr := s.db.MarkRunAsCompleted(ctx, db.MarkRunAsCompletedParams{
				ID:                  runID,
				FilesProcessedCount: pgtype.Int4{Int32: int32(numFiles), Valid: true},
			}); statusErr != nil {
				ctxlog.Error("Failed to update job status", "error", statusErr)
			}

			ctxlog.Info("Successfully indexed repository", "num_files", numFiles)
		}()

		ctxlog.Info("Started indexing repository")
	}
}

func (s *Server) ingestFiles(
	ctx context.Context,
	qs *db.Queries,
	runID, slug, sha string,
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
