package github

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login     string `json:"login"`
		URL       string `json:"html_url"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
}
