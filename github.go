package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ErrRateLimit struct {
	Limit     int
	Remaining int
	Used      int
	Reset     time.Time
}

func (e ErrRateLimit) Error() string {
	return fmt.Sprintf("rate limit exceeded")
}

type GitHubClient struct {
	httpClient *http.Client
	apiToken   string
}

func newGitHubClient(apiToken string) *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{},
		apiToken:   apiToken,
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
	URL           string `json:"html_url"`
	Description   string `json:"description"`
	DefaultBranch string `json:"default_branch"`
}

func newErrRateLimit(header http.Header) ErrRateLimit {
	var (
		remaining, _ = strconv.Atoi(header.Get("X-RateLimit-Remaining"))
		limit, _     = strconv.Atoi(header.Get("X-RateLimit-Limit"))
		reset, _     = strconv.Atoi(header.Get("X-RateLimit-Reset"))
		used, _      = strconv.Atoi(header.Get("X-RateLimit-Limit"))
	)
	return ErrRateLimit{
		Limit:     limit,
		Remaining: remaining,
		Used:      used,
		Reset:     time.Unix(int64(reset), 0),
	}
}

func (c *GitHubClient) GetRepository(
	ctx context.Context,
	owner, name string,
) (*Repository, []byte, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return nil, nil, newErrRateLimit(resp.Header)
		}
		return nil, nil, fmt.Errorf("failed to get repository: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var repo Repository
	if err := json.Unmarshal(b, &repo); err != nil {
		return nil, nil, err
	}

	return &repo, b, nil
}

func (c *GitHubClient) GetLatestCommitSHA(
	ctx context.Context,
	owner, name, branch string,
) (string, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, name, branch)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return "", newErrRateLimit(resp.Header)
		}
		return "", fmt.Errorf("failed to get repository: %s", resp.Status)
	}

	var response struct {
		SHA string `json:"sha"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.SHA, nil
}
