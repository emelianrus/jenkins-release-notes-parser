package web

/*

we have from github
  all versions release notes

  1) check if we have plugin data in cache?
  if (exist) {
	use
  } else {
	get release from github
	add to cache
  }

*/
import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/github"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
)

// HTML start
// part of html response
type Version struct {
	Version string
	Changes template.HTML
}

// represent repo in github
type Product struct {
	Name             string
	Versions         []Version
	InstalledVersion string
	LastUpdated      string // TODO
}

// HTML end
var ownerName = "jenkinsci"

// get jenkins server (hardcoded)
// get plugins from jenikins server
// check cache for plugin by versions file
// construct pageData
func getReleaseNotesPageData(redisclient *db.Redis, jenkinsServer types.JenkinsServer) ([]Product, error) {
	// default page data

	products := []Product{}

	for _, plugin := range jenkinsServer.Plugins {

		pluginVersionsJson, err := redisclient.GetPluginVersions(plugin.Name)
		if err != nil {
			fmt.Println("versions file doesn't exist in redis cache for " + plugin.Name)
			fmt.Println(err)
			releases, err := github.Download(plugin.Name)
			if err != nil {
				fmt.Println("Failed to get releases from github")
				continue
			}
			err = redisclient.SaveReleaseNotesToDB(releases, plugin.Name)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Failed to save release notes to db")
			}

			pluginVersionsJson, err = redisclient.GetPluginVersions(plugin.Name)
			if err != nil {
				fmt.Println(err)
				fmt.Println("2nd attempt to GetPluginVersions failed")
				// return web page with default values
				return []Product{}, errors.New("2nd attempt to GetPluginVersions failed")
			}
		}

		// Assume we hit redis cache
		var versions []string
		err = json.Unmarshal(pluginVersionsJson, &versions)
		if err != nil {
			log.Println(err)
			// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
			return []Product{}, errors.New("Failed to unmarshal releases from cache")
		}

		var convertedVersions []Version
		// TODO: check jenkins plugin version and show only diff from installed version to latest
		for _, version := range versions {

			releaseNote, err := redisclient.GetPluginWithVersion(plugin.Name, version)
			if err != nil {
				log.Println(err)
				// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
				return []Product{}, errors.New("Failed to unmarshal releases notes from cache")
			}

			convertedVersions = append(convertedVersions, Version{
				Version: version,
				Changes: template.HTML(
					utils.ReplaceGitHubLinks(
						utils.ConvertMarkDownToHtml(releaseNote.Body))),
			})
		}

		lastUpdated, _ := redisclient.Get(fmt.Sprintf("github:%s:%s:%s", ownerName, plugin.Name, "lastUpdated")).Bytes()

		products = append(products,
			Product{
				Name:             plugin.Name,
				Versions:         convertedVersions,
				InstalledVersion: plugin.Version,
				LastUpdated:      string(lastUpdated),
			},
		)
	}
	return products, nil
}
