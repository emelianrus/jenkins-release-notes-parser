package types

type GitHubReleases struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Notes   string `json:"body"`
	Url     string `json:"url"`
}
