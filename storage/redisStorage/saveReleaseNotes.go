package redisStorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

func (r *RedisStorage) SaveReleaseNotesToDB(releases []types.ReleaseNote, projectName string) error {
	currentTime := time.Now()
	formattedTime := currentTime.Format("02 January 2006 15:04:05")

	err := r.SetLastUpdatedTime(projectName, formattedTime)
	if err != nil {
		log.Println(err)
		return errors.New("set lastUpdated failed")
	}

	var versions []string

	// save repo release notes per version
	for _, release := range releases {

		// TODO: some plugins doesnt have name, so replace with tag
		if release.Name == "" {
			release.Name = release.Tag
		}

		versions = append(versions, release.Name)

		jsonData, err := json.Marshal(release)
		if err != nil {
			// log.Println(err)
			return fmt.Errorf("error Marshal release: %s", err)
		}
		// save release file
		key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, release.Name)
		err = r.DB.Set(key, jsonData)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error setting release: %s", err)
		}
	}

	// save "versions" file
	jsonVersions, _ := json.Marshal(versions)
	err = r.DB.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "versions"),
		jsonVersions)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}

	if len(versions) == 0 {
		fmt.Println("Project doesn't have releases: " + projectName)
		r.SetProjectError(projectName, "project doesn't have releases")
		return fmt.Errorf("project doesn't have releases")
	}
	// save "latestVersion" file
	jsonLatestVersion, _ := json.Marshal(versions[0])
	err = r.DB.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "latestVersion"),
		jsonLatestVersion)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}

	return nil
}