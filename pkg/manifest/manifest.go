package manifest

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

// code to read jenkins/java MANIFEST.MF file and parse it

// Represent MANIFEST.MF java file, holds map of attributes
type Manifest map[string]string

// Represent dependency of plugin
type Dependency struct {
	Name       string
	Version    string
	Optional   bool
	Resolution bool
}

// ReadFile reads JAR file and parses manifest file
func Parse(hpiFile string) (Manifest, error) {
	logrus.Debugf("ReadMainfest file: %s", hpiFile)

	r, err := zip.OpenReader(hpiFile)

	if err != nil {
		return nil, err
	}

	defer r.Close()
	for _, f := range r.File {
		if f.Name != "META-INF/MANIFEST.MF" {
			continue
		}

		rc, err := f.Open()

		if err != nil {
			return nil, err
		}
		return parseManifestData(rc)
	}

	return nil, errors.New("given file is not a HPI file")
}

func (mf Manifest) GetDependencies() []Dependency {
	var deps []Dependency

	if len(mf["Plugin-Dependencies"]) < 1 {
		return nil
	}
	lines := strings.Split(mf["Plugin-Dependencies"], ",")
	for _, line := range lines {
		var dep Dependency

		if strings.Contains(line, ";resolution:=optional") {
			dep.Optional = true
		} else {
			dep.Optional = false
		}
		line = strings.ReplaceAll(line, ";resolution:=optional", "")

		plugin := strings.Split(line, ":")
		dep.Name = plugin[0]
		dep.Version = plugin[1]

		deps = append(deps, dep)
	}

	if len(deps) < 1 {
		fmt.Println("ERROR DEPS IS NIL")
	}
	return deps
}

func parseManifestData(r io.Reader) (Manifest, error) {
	m := make(Manifest)
	s := bufio.NewScanner(r)

	var propName, propVal string
	for s.Scan() {
		text := s.Text()

		if len(text) == 0 {
			continue
		}

		if strings.HasPrefix(text, " ") {
			m[propName] += strings.TrimLeft(text, " ")
			continue
		}

		propSepIndex := strings.Index(text, ": ")

		if propSepIndex == -1 || len(text) < propSepIndex+2 {
			return nil, errors.New("can't parse manifest file (wrong format)")
		}

		propName = text[:propSepIndex]
		propVal = text[propSepIndex+2:]
		m[propName] = propVal
	}
	if len(m) < 1 {
		fmt.Println("error reading manifest")
		return nil, errors.New("error reading manifest, its nil")
	}
	return m, nil
}
