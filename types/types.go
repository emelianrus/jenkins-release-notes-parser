package types

// Represent one release/release note in project
type ReleaseNote struct {
	Name      string
	Tag       string
	BodyHTML  string
	HTMLURL   string
	CreatedAt string
}

type Project struct {
	Name  string
	Owner string

	Error        string
	IsDownloaded bool
	LastUpdated  string

	// TODO: should be here
	ReleaseNotes []ReleaseNote
}
