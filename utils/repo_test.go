package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepoRef_URL(t *testing.T) {
	repo := RepoRef{
		Provider: "github.com",
		Owner:    "marcusolsson",
		RepoName: "recipes",
	}

	require.Equal(t, "https://github.com/marcusolsson/recipes", repo.URL().String())
	require.Equal(t, "github.com/marcusolsson/recipes", repo.ID())
}

func TestRepoFileRef_URL(t *testing.T) {
	repo := RepoRef{
		Provider: "github.com",
		Owner:    "marcusolsson",
		RepoName: "recipes",
	}

	fileRef := RepoFileRef{
		Repo: repo,
		Path: "Popcorn.cook",
		Ref:  "main",
	}

	require.Equal(
		t,
		"https://github.com/marcusolsson/recipes/blob/main/Popcorn.cook",
		fileRef.URL().String(),
	)

	require.Equal(
		t,
		"github.com/marcusolsson/recipes/Popcorn.cook@main",
		fileRef.ID(),
	)
}
