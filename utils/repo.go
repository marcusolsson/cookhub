package utils

import "net/url"

type RepoRef struct {
	Provider string
	Owner    string
	RepoName string
}

func (r RepoRef) URL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   r.Provider,
		Path:   "/" + r.Owner + "/" + r.RepoName,
	}
}

func (r RepoRef) ID() string {
	return r.Provider + "/" + r.Owner + "/" + r.RepoName
}

func (r RepoRef) Slug() string {
	return r.Owner + "/" + r.RepoName
}

type RepoFileRef struct {
	Repo RepoRef
	Path string
	Ref  string
}

func (r RepoFileRef) ID() string {
	return r.Repo.ID() + "/" + r.Path + "@" + r.Ref
}

func (r RepoFileRef) URL() *url.URL {
	repoURL := r.Repo.URL()

	u := &url.URL{
		Path: repoURL.Path + "/blob/" + r.Ref + "/" + r.Path,
	}

	return r.Repo.URL().ResolveReference(u)
}
