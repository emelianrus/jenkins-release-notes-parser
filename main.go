package main

import (
	"fmt"
	"io/ioutil"
	"os"

	git "github.com/emelianrus/jenkins-release-notes-parser/pkg/releasenotes/git"
	releaseParser "github.com/emelianrus/jenkins-release-notes-parser/pkg/releasenotes/releaseParser"
)

func main() {
	byt, _ := ioutil.ReadFile("diff.out")
	diff, _ := git.Parse(string(byt))

	project, _ := releaseParser.GetChangedVersions(diff)
	for _, f := range project.Files {

		var releaseNotes string
		for _, p := range f.Plugins {
			releases, err := releaseParser.GetReleases(p.Name)
			if err != nil {
				fmt.Println(err)
				break
			}
			_, result := releaseParser.GetReleaseNotes(releases, p)

			releaseNotes += result
		}

		err := os.WriteFile(f.Name+"_RELEASE_NOTES.md", []byte(releaseNotes), 0644)

		if err != nil {
			fmt.Println(err)
		}

	}
}
