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

func downloadZipBall(ctx context.Context, repoSlug, commitSHA string) (string, error) {
	folderPath, err := os.MkdirTemp("", "cooklang-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(folderPath)

	u := fmt.Sprintf("https://github.com/%s/zipball/%s", repoSlug, commitSHA)

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

type File struct {
	Name    string
	Content []byte
}

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
					relPath := strings.TrimPrefix(file.Name, archiveName)

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
