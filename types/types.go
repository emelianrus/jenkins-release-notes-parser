package types

// response from github getted from release path
type GitHubReleaseNote struct {
	Name      string `json:"name"` // Version
	Body      string `json:"body"` // this is markdown formated text of release note
	CreatedAt string `json:"created_at"`
}

type JenkinsPlugin struct {
	Name         string
	Version      string
	Error        string
	IsDownloaded bool
}

// watch point like group with watched project
type JenkinsServer struct {
	Name    string
	Core    string
	Plugins []JenkinsPlugin
}
