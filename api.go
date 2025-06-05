package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	db "github.com/marcusolsson/cookhub/db/sqlc"
)

func (s *server) apiIndexRepo(w http.ResponseWriter, req *http.Request) {
	var (
		org      = chi.URLParam(req, "org")
		repoName = chi.URLParam(req, "name")
		slug     = fmt.Sprintf("%s/%s", org, repoName)
	)

	ctx := req.Context()

	ctxlog := s.logger.With("slug", slug)

	sha, err := s.ghClient.GetLatestCommitSHA(ctx, slug)
	if err != nil {
		ctxlog.Error("Failed to get latest commit SHA from repo", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.runJob(ctx, slug, sha); err != nil {
		ctxlog.Error("Job failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) runJob(ctx context.Context, slug, sha string) (jobErr error) {
	jobID, err := s.db.CreateJob(ctx, db.CreateJobParams{Slug: slug, CommitSha: sha})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			s.db.SetJobStatus(ctx, db.SetJobStatusParams{
				ID:     jobID,
				Status: db.StatusEnumFailed,
			})
		} else {
			s.db.SetJobStatus(ctx, db.SetJobStatusParams{
				ID:     jobID,
				Status: db.StatusEnumCompleted,
			})
		}
	}()

	tempDir, err := os.MkdirTemp("", "cooklang-")
	if err != nil {
		jobErr = err
		s.logger.Error("Failed to create temporary folder", "error", err)
		return
	}
	defer os.RemoveAll(tempDir)

	zipPath, err := downloadZipBall(ctx, slug, sha, tempDir)
	if err != nil {
		jobErr = err
		s.logger.Error("Failed to download zipball", "error", err)
		return
	}

	var params []db.CreateFileParams
	for file, err := range readFilesFromZip(zipPath) {
		if err != nil {
			jobErr = err
			s.logger.Error("Failed to read file from zip", "error", err)
			return
		}

		params = append(params, db.CreateFileParams{
			JobID:   jobID,
			Name:    file.Name,
			Content: string(file.Content),
		})
	}

	_, err = s.db.CreateFile(ctx, params)
	if err != nil {
		jobErr = err
		s.logger.Error("Failed to create files in database", "error", err)
		return
	}

	s.logger.Info(
		"Successfully indexed repository",
		"slug", slug,
		"commit_sha", sha,
		"num_files", len(params),
	)

	return
}
