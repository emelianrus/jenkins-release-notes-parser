package jenkins

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/sirupsen/logrus"
)

type PluginSite struct{}

func NewPluginSite() PluginSite {
	return PluginSite{}
}

type pluginSiteRelease struct {
	TagName     string `json:"tagName"`
	Name        string `json:"name"`
	PublishedAt string `json:"publishedAt"`
	HTMLURL     string `json:"htmlURL"`
	BodyHTML    string `json:"bodyHTML"`
}

type PluginSiteReleases struct {
	Releases []pluginSiteRelease `json:"releases"`
}

func (ps *PluginSite) Download(projectName string) ([]types.ReleaseNote, error) {
	logrus.Infoln("[PluginSiteReleases][Download]")
	releaseNotes := []types.ReleaseNote{}

	url := fmt.Sprintf("https://plugin-site-issues.jenkins.io/api/plugin/%s/releases", projectName)
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorln("Error creating request:", err)
		return nil, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorln("Error making request:", err)
		return nil, errors.New("Error making request")
	}

	var releases PluginSiteReleases
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		logrus.Errorln("can not decode plugin site releases")
		return nil, nil
	}

	for _, release := range releases.Releases {
		releaseNotes = append(releaseNotes, types.ReleaseNote{
			Name:      release.Name,
			Tag:       release.TagName,
			BodyHTML:  release.BodyHTML,
			HTMLURL:   release.HTMLURL,
			CreatedAt: release.PublishedAt,
		})
	}

	return releaseNotes, nil
}
