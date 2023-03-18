package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

func (r *Redis) SaveReleaseNotesToDB(releases []types.GitHubReleaseNote, pluginName string) error {

	currentTime := time.Now()
	formattedTime := currentTime.Format("02 January 2006 15:04")

	// set lastUpdated file for repo
	err := r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "lastUpdated"),
		formattedTime)
	if err != nil {
		log.Println(err)
		return errors.New("set lastUpdated failed")
	}

	var versions []string

	// save repo release notes per version
	for _, release := range releases {
		versions = append(versions, release.Name)
		key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, release.Name)

		jsonData, err := json.Marshal(release)
		if err != nil {
			// log.Println(err)
			return fmt.Errorf("error Marshal release: %s", err)
		}
		err = r.Set(key, jsonData)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error setting release: %s", err)
		}
	}

	// save "versions" file
	jsonVersions, _ := json.Marshal(versions)
	err = r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "versions"),
		jsonVersions)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}

	// save "latestVersion" file
	jsonLatestVersion, _ := json.Marshal(versions[0])
	err = r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "latestVersion"),
		jsonLatestVersion)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}

	return nil
}
