package releaseParser

import (
	"fmt"

	"regexp"
	"strings"

	types "github.com/emelianrus/jenkins-release-notes-parser/types"
)

type ReleaseNotes struct {
	Name    string
	TagName string
	Notes   string
}

// jenkins introduced new version style, we should support both
func normalizeVersions(tagName string, pluginName string) string {
	fmt.Println("checking" + pluginName + ":" + tagName)
	match, _ := regexp.MatchString(".*\\.v.*", tagName)

	if match {
		return tagName
	} else {
		if len(strings.Split(tagName, pluginName+"-")) == 2 {
			return strings.Split(tagName, pluginName+"-")[1]
		} else {
			return tagName
		}
	}
}

// function returns notes between two versions
func GetReleaseNotes(releases []types.GitHubReleases, plugin types.Plugin) (bool, string) {

	var upgrade bool
	var result string

	var NewVersionFound bool
	var OldVersionFound bool

	if len(releases) < 1 {
		return true, "# " + plugin.Name + ":" + plugin.NewVersion + "\n" + "RELEASES NOT FOUND" + "\n --- \n"
	}

	for _, release := range releases {
		tagName := normalizeVersions(release.TagName, plugin.Name)

		// to detect its upgrade or downgrade of plugin
		if tagName == plugin.NewVersion {
			NewVersionFound = true
			if !OldVersionFound {
				upgrade = true
			}
		}

		if tagName == plugin.OldVersion {
			OldVersionFound = true
			if !NewVersionFound {
				upgrade = false
			}
		}

		if NewVersionFound || OldVersionFound {
			if upgrade && OldVersionFound {
				break
			}
			result += "# " + plugin.Name + ":" + release.TagName + "\n" + release.Notes + "\n --- \n"

			if NewVersionFound && OldVersionFound {
				break
			}
		}
	}

	return upgrade, result
}
