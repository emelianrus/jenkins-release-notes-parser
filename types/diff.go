package types

type FileMode int

const (
	DELETED FileMode = iota
	MODIFIED
	NEW
)

func (d FileMode) String() string {
	return [...]string{"DELETED", "MODIFIED", "NEW"}[d]
}

type DiffLineMode rune

const (
	ADDED DiffLineMode = iota
	REMOVED
	UNCHANGED
)

func (d DiffLineMode) String() string {
	return [...]string{"ADDED", "REMOVED", "UNCHANGED"}[d]
}

// DiffLine is the least part of an actual diff
type DiffLine struct {
	Mode    DiffLineMode
	Content string
}

// DiffChunk is a group of difflines
type DiffChunk struct {
	Lines []*DiffLine
}

// DiffFile is the sum of DiffChunks and holds the changes of the file features
type DiffFile struct {
	Mode     FileMode
	OrigName string
	NewName  string
	Chunks   []*DiffChunk
}

// Diff is the collection of DiffFiles
type Diff struct {
	Files []*DiffFile
	Raw   string `sql:"type:text"`

	PullID uint `sql:"index"`
}

// plugin
type Plugin struct {
	Name       string
	NewVersion string
	OldVersion string
}

type File struct {
	Name    string
	Plugins []Plugin
}

type Project struct {
	Files []*File
}
