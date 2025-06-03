package api

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mattn/go-sqlite3"
)

func (s *APIServer) getIngestedFiles(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var (
		org  = chi.URLParam(req, "org")
		name = chi.URLParam(req, "name")
		slug = fmt.Sprintf("%s/%s", org, name)
	)

	jobID, err := s.Store.GetLatestJobBySlug(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := s.Store.GetFilesByJob(ctx, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *APIServer) indexRepo(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var (
		org      = chi.URLParam(req, "org")
		repoName = chi.URLParam(req, "name")
		slug     = fmt.Sprintf("%s/%s", org, repoName)
	)

	tempDir, err := os.MkdirTemp("staging", "cooklang-")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// metadata, err := getRepoMetadata(ctx, slug)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	zipPath, err := downloadZipBall(ctx, slug, tempDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath))

	sha := strings.Split(name, "-")[2]

	jobID, err := s.Store.CreateJob(ctx, slug, sha)
	if err != nil {
		var sqlite3Err sqlite3.Error

		if errors.As(err, &sqlite3Err) && sqlite3Err.Code == sqlite3.ErrConstraint {
			fmt.Println("Already ingested")
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if filepath.Ext(file.Name) == ".cook" {
			fi := file.FileInfo()

			if !fi.IsDir() {
				b, err := readZipFile(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				archiveName := strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath))
				relPath := strings.TrimPrefix(file.Name, archiveName)

				if err := s.Store.CreateFile(ctx, jobID, relPath, b); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

// type Repository struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	FullName string `json:"full_name"`
// 	Owner    struct {
// 		Login string `json:"login"`
// 		ID    int    `json:"id"`
// 		URL   string `json:"html_url"`
// 	}
// 	URL             string `json:"html_url"`
// 	Description     string `json:"description"`
// 	StargazersCount int    `json:"stargazers_count"`
// 	WatchersCount   int    `json:"watchers_count"`
// }
//
// func getRepoMetadata(ctx context.Context, repoSlug string) ([]byte, error) {
// 	u := fmt.Sprintf("https://api.github.com/repos/%s", repoSlug)
//
// 	resp, err := http.Get(u)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
//
// 	return io.ReadAll(resp.Body)
// }

func downloadZipBall(ctx context.Context, repoSlug string, folderPath string) (string, error) {
	u := fmt.Sprintf("https://github.com/%s/zipball/HEAD", repoSlug)

	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dispo := resp.Header.Get("Content-Disposition")

	mediaType, params, err := mime.ParseMediaType(dispo)
	if err != nil {
		return "", err
	}

	if mediaType != "attachment" {
		return "", fmt.Errorf("unexpected media type: %s", mediaType)
	}

	zipPath := filepath.Join(folderPath, params["filename"])

	f, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	return zipPath, nil
}

func readZipFile(f *zip.File) ([]byte, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}
