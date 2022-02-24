package git

import (
	"errors"
	"strings"

	types "github.com/emelianrus/jenkins-release-notes-parser/types"
)

func lineMode(line string) (*types.DiffLineMode, error) {
	var mode types.DiffLineMode
	switch line[:1] {
	case " ":
		mode = types.UNCHANGED
	case "+":
		mode = types.ADDED
	case "-":
		mode = types.REMOVED
	default:
		return nil, errors.New("could not parse line mode for line: \"" + line + "\"")
	}
	return &mode, nil
}

// Parse takes a diff, such as produced by "git diff", and parses it into a
// Diff struct.
func Parse(diffString string) (*types.Diff, error) {
	var diff types.Diff
	diff.Raw = diffString
	lines := strings.Split(diffString, "\n")

	var file *types.DiffFile
	var chunk *types.DiffChunk
	var isChunk bool
	var isNewLineSymbol bool
	isNewLineSymbol = false
	// Parse each line of diff.
	for _, line := range lines {
		switch {

		case isNewLineSymbol:
			isNewLineSymbol = false
			continue
		case line == `\ No newline at end of file`:
			isNewLineSymbol = true
			chunk.Lines = chunk.Lines[:len(chunk.Lines)-1]
			continue
		case strings.HasPrefix(line, "diff "):
			isChunk = false
			// Start a new file.
			file = &types.DiffFile{}
			diff.Files = append(diff.Files, file)
			// File mode.
			file.Mode = types.MODIFIED
		case line == "+++ /dev/null":
			file.Mode = types.DELETED
		case line == "--- /dev/null":
			file.Mode = types.NEW
		case strings.HasPrefix(line, "--- a/"):
			file.OrigName = strings.TrimPrefix(line, "--- a/")
		case strings.HasPrefix(line, "--- a/"):
			file.NewName = strings.TrimPrefix(line, "--- a/")
		case strings.HasPrefix(line, "@@ "):
			isChunk = true
			// Start new chunk.
			chunk = &types.DiffChunk{}
			file.Chunks = append(file.Chunks, chunk)

		case isChunk && isSourceLine(line):
			mode, err := lineMode(line)
			if err != nil {
				return nil, err
			}

			line := &types.DiffLine{
				Mode:    *mode,
				Content: line[1:],
			}
			chunk.Lines = append(chunk.Lines, line)
		}
	}

	return &diff, nil
}

func isSourceLine(line string) bool {
	if line == `\ No newline at end of file` {
		return false
	}
	if l := len(line); l == 0 || (l >= 3 && (line[:3] == "---" || line[:3] == "+++")) {
		return false
	}
	return true
}
