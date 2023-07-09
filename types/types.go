package types

// Represent one release/release note in project
type ReleaseNote struct {
	Name      string
	Tag       string
	BodyHTML  string
	HTMLURL   string
	CreatedAt string
}
