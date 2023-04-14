package types

// main type of downloaded project
type ReleaseNote struct {
	Name      string
	Tag       string
	BodyHTML  string
	CreatedAt string
}

type Project struct {
	Name         string
	Owner        string
	Version      string
	Error        string
	IsDownloaded bool
	LastUpdated  string
}

// watch point like group with watched project
type JenkinsServer struct {
	Name    string
	Core    string
	Plugins []Project
}
