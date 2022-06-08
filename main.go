package main

// "fmt"
// "io/ioutil"
// "os"
import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/releasenotes/git"
	releaseParser "github.com/emelianrus/jenkins-release-notes-parser/pkg/releasenotes/releaseParser"
)

func main() {
	byt, _ := ioutil.ReadFile("diff.out")
	diff, _ := git.Parse(string(byt))
	fmt.Println(diff)
	project, _ := releaseParser.GetChangedVersions(diff)

	for _, localPlugins := range project.Files {

		var releaseNotes string
		for _, lp := range localPlugins.Plugins {
			releases, err := releaseParser.GetPluginReleases(lp.Name)
			if err != nil {
				fmt.Println(err)
				break
			}

			parasedReleases := releases.GetReleasesBetweenVersions(lp.OldVersion, lp.NewVersion)

			for _, rel := range parasedReleases {
				releaseNotes += "# " + rel.Name + ":" + rel.TagName + "\n" + rel.Notes + "\n --- \n"

			}
		}

		err := os.WriteFile(localPlugins.Name+"_RELEASE_NOTES.md", []byte(releaseNotes), 0644)

		if err != nil {
			fmt.Println(err)
		}

	}

}
