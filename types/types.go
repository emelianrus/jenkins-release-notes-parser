package types

type GitHubReleaseNote struct {
	Name      string `json:"name"` // Version
	Body      string `json:"body"` // this is markdown formated text of release note
	CreatedAt string `json:"created_at"`
}

type JenkinsPlugin struct {
	Name    string
	Version string
}

type JenkinsServer struct {
	Name    string
	Core    string
	Plugins []JenkinsPlugin
}

// From redis
type AllVersions []string
