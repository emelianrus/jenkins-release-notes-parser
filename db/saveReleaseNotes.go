package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/sirupsen/logrus"
)

func (r *Redis) SaveReleaseNotesToDB(releases []types.GitHubReleaseNote, projectName string) error {

	currentTime := time.Now()
	formattedTime := currentTime.Format("02 January 2006 15:04")

	logrus.Debugf("update latestUpdated time to: %v", currentTime)

	// set lastUpdated file for repo
	err := r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "lastUpdated"),
		formattedTime)
	if err != nil {
		log.Println(err)
		return errors.New("set lastUpdated failed")
	}

	var versions []string

	// save repo release notes per version
	for _, release := range releases {

		// TODO: some plugins doesnt have name, so replace with tag
		if release.Name == "" {
			release.Name = release.TagName
		}
		versions = append(versions, release.Name)
		key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, release.Name)

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
	err = r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "versions"),
		jsonVersions)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}
	if len(versions) == 0 {
		fmt.Println("Project doesn't have releases: " + projectName)
		return nil
	}
	// save "latestVersion" file
	jsonLatestVersion, _ := json.Marshal(versions[0])
	err = r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "latestVersion"),
		jsonLatestVersion)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}

	return nil
}
