package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GitHubClient struct {
	httpClient *http.Client
}

func newGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{},
	}
}

type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
		URL   string `json:"html_url"`
	}
	URL             string `json:"html_url"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
}

func (c *GitHubClient) GetRepository(ctx context.Context, repoSlug string) (*Repository, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s", repoSlug)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (c *GitHubClient) GetLatestCommitSHA(ctx context.Context, repoSlug string) (string, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/commits/HEAD", repoSlug)

	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response struct {
		SHA string `json:"sha"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.SHA, nil
}
