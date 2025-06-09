package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"iter"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// downloadZipBall downloads a zip archive of a specific commit from a GitHub
// repository. The archive is saved to the specified directory and the path to
// the zip file is returned.
func downloadZipBall(ctx context.Context, repoSlug, commitSHA, dirPath string) (string, error) {
	u := fmt.Sprintf("https://github.com/%s/zipball/%s", repoSlug, commitSHA)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}

	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	disposition := resp.Header.Get("Content-Disposition")

	mediaType, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		return "", err
	}

	if mediaType != "attachment" {
		return "", fmt.Errorf("unexpected media type: %s", mediaType)
	}

	zipPath := filepath.Join(dirPath, params["filename"])

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

type File struct {
	Name    string
	Content []byte
}

// readFilesFromZip reads .cook files from a zip archive and yields them one by one.
func readFilesFromZip(zipPath string) iter.Seq2[*File, error] {
	return func(yield func(*File, error) bool) {
		reader, err := zip.OpenReader(zipPath)
		if err != nil {
			yield(nil, err)
			return
		}
		defer reader.Close()

		for _, file := range reader.File {
			if filepath.Ext(file.Name) == ".cook" {
				fi := file.FileInfo()

				if !fi.IsDir() {
					b, err := readZipFile(file)
					if err != nil {
						if !yield(nil, err) {
							return
						}
					}

					archiveName := strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath))
					relPath := strings.TrimPrefix(file.Name, archiveName+"/")

					if !yield(&File{
						Name:    relPath,
						Content: b,
					}, nil) {
						return
					}
				}
			}
		}
	}
}
