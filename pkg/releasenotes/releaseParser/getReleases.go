package releaseParser

import (
	"fmt"
	"log"

	"encoding/json"
	"errors"
	"strings"

	types "github.com/emelianrus/jenkins-release-notes-parser/types"
)

// curl \
//   -H "Accept: application/vnd.github.v3+json" \
//   https://api.github.com/repos/jenkinsci/configuration-as-code-plugin/releases
// TODO: https://docs.github.com/en/rest/reference/rate-limit
const GITHUB_API = "https://api.github.com/repos/"

func GetReleases(pluginName string) ([]types.GitHubReleases, error) {
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
	var releases []types.GitHubReleases

	err := json.Unmarshal(body, &releases)
	if err != nil {
		fmt.Println(string(body))
		// TODO: fix error handling to write
		// log.Fatalln(err)
		log.Println(err)
		return []types.GitHubReleases{}, errors.New("error to get URL(rate limit) or unmarhal json")
	}

	return releases, err
}
