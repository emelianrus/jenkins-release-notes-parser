package main

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
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/utils"
)

type GitHubReleaseNote struct {
	Name      string `json:"name"` // Version
	Body      string `json:"body"` // this is markdown formated text of release note
	CreatedAt string `json:"created_at"`
}

// From redis
type Versions []string

// HTML start
// part of html responce
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

func saveReleaseNotesToDB(redisclient *Redis, releases []GitHubReleaseNote, pluginName string) error {

	currentTime := time.Now()
	formattedTime := currentTime.Format("02 January 2006 15:04")
	err := redisclient.SetLastUpdatedTime(pluginName, formattedTime)
	if err != nil {
		// fmt.Println(err)
		// fmt.Println("Can not set updated time")
		return fmt.Errorf("error setting lastUpdate time in get github release: %s", err)
	}
	versions := Versions{}

	for _, release := range releases {
		versions = append(versions, release.Name)
		key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, release.Name)
		// 0 time.Hour
		jsonData, err := json.Marshal(release)
		if err != nil {
			// log.Println(err)
			return fmt.Errorf("error Marshal release: %s", err)
		}
		err = redisclient.Set(key, jsonData)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error setting release: %s", err)
		}
	}

	jsonVersions, _ := json.Marshal(versions)
	err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "versions"),
		jsonVersions)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}
	return nil
}

// get jenkins server (hardcoded)
// get plugins from jenikins server
// check cache for plugin by versions file
// construct pageData
func getReleaseNotesPageData(redisclient *Redis, jenkinsServer JenkinsServer) ([]Product, error) {
	// default page data

	products := []Product{}

	for _, plugin := range jenkinsServer.Plugins {

		pluginVersionsJson, err := redisclient.GetPluginVersions(plugin.Name)
		if err != nil {
			fmt.Println("versions file doesn't exist in redis cache for " + plugin.Name)
			fmt.Println(err)
			releases, err := GetGitHubReleases(plugin.Name)
			if err != nil {
				fmt.Println("Failed to get releases from github")
				continue
			}
			err = saveReleaseNotesToDB(redisclient, releases, plugin.Name)
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
		var versions Versions
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
				Changes: template.HTML(utils.ReplaceGitHubLinks(utils.ConvertMarkDownToHtml(releaseNote.Body))),
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
