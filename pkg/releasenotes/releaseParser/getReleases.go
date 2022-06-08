package releaseParser

import (
	"fmt"
	"log"
	"regexp"

	"encoding/json"
	"errors"
	"strings"
)

// curl \
//   -H "Accept: application/vnd.github.v3+json" \
//   https://api.github.com/repos/jenkinsci/configuration-as-code-plugin/releases
// TODO: https://docs.github.com/en/rest/reference/rate-limit
const GITHUB_API = "https://api.github.com/repos/"

type GitHubRelease struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Notes   string `json:"body"`
	Url     string `json:"url"`
}

type GitHubReleases struct {
	GitHubReleases []GitHubRelease
}

func GetPluginReleases(pluginName string) (GitHubReleases, error) {
	// Fixes workarounds :)
	// plugins refused go by standards
	postFix := "-plugin"
	if strings.Contains(pluginName, "-plugin") {
		postFix = ""
	}

	if pluginName == "cloudbees-bitbucket-branch-source" {
		pluginName = "bitbucket-branch-source"
	}

	URL := GITHUB_API + "jenkinsci/" + pluginName + postFix + "/releases"
	body := SendRequest(URL)
	var releases []GitHubRelease

	err := json.Unmarshal(body, &releases)

	if err != nil {
		fmt.Println(string(body))
		// TODO: fix error handling to write
		// log.Fatalln(err)
		log.Println(err)
		return GitHubReleases{GitHubReleases: releases}, errors.New("error to get URL(rate limit) or unmarhal json")
	}

	return GitHubReleases{GitHubReleases: releases}, err
}

// HELPER
// jenkins introduced new version style, we should support both
func normalizeVersions(tagName string, pluginName string) string {
	fmt.Println("checking " + pluginName + ":" + tagName)
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

func (gr GitHubReleases) PrintVersions() {
	for _, release := range gr.GitHubReleases {
		fmt.Println(release.TagName)
	}
}

func (gr GitHubReleases) GetLatestRelease() GitHubRelease {
	return gr.GitHubReleases[0]
}

// returns releases between currently installed version and new(changed git)
func (gr GitHubReleases) GetReleasesBetweenVersions(oldVersion string, newVersion string) []GitHubRelease {
	fmt.Println("versions between old: " + oldVersion + " new: " + newVersion)
	var resultReleases []GitHubRelease

	var upgrade bool
	var result string

	var NewVersionFound bool
	var OldVersionFound bool

	// if len(gr.GitHubReleases) < 1 {
	// 	return true, "# " + plugin.Name + ":" + plugin.NewVersion + "\n" + "RELEASES NOT FOUND" + "\n --- \n"
	// }

	for _, release := range gr.GitHubReleases {

		tagName := normalizeVersions(release.TagName, release.Name)
		fmt.Println("tagname: " + tagName)
		// to detect its upgrade or downgrade of plugin
		if tagName == newVersion {
			NewVersionFound = true
			if !OldVersionFound {
				upgrade = true
			}
		}

		if tagName == oldVersion {
			OldVersionFound = true
			if !NewVersionFound {
				upgrade = false
			}
		}

		if NewVersionFound || OldVersionFound {
			if upgrade && OldVersionFound {
				break
			}
			resultReleases = append(resultReleases, release)
			result += "# " + release.Name + ":" + release.TagName + "\n" + release.Notes + "\n --- \n"

			if NewVersionFound && OldVersionFound {
				break
			}
		}
	}

	return resultReleases
}
