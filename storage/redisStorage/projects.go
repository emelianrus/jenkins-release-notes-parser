package redisStorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
	"github.com/sirupsen/logrus"
)

func (r *RedisStorage) GetAllProjects() ([]types.Project, error) {
	var projects []types.Project
	// TODO: this should be configurable
	repoOwner := "jenkinsci"

	// get list of all projects
	var projectsTmp []string

	projectsKeys, _ := r.DB.Keys(fmt.Sprintf("github:%s:*", repoOwner))

	for _, key := range projectsKeys {
		splitted := strings.Split(key, ":")
		projectsTmp = append(projectsTmp, splitted[2])
	}
	projectsTmp = utils.RemoveDuplicates(projectsTmp)

	// gather project data
	for _, projectName := range projectsTmp {

		projects = append(projects, types.Project{
			Name:         projectName,
			Owner:        repoOwner,
			Error:        r.GetProjectError(repoOwner, projectName),
			IsDownloaded: r.IsProjectDownloaded(repoOwner, projectName),
			LastUpdated:  r.GetLastUpdatedTime(repoOwner, projectName),
		})

	}

	return projects, nil
}

// get one project with release notes
func (r *RedisStorage) GetProjectReleaseNotes(projectOwner string, projectName string) ([]types.ReleaseNote, error) {

	releaseNotes := []types.ReleaseNote{}

	// get all versions for specific project
	projectVersionsJson, err := r.DB.Get(fmt.Sprintf("github:%s:%s:versions", projectOwner, projectName))
	if err != nil {
		return releaseNotes, errors.New("error in getPlugins " + projectName)
	}
	// convert json versions to []string
	var versions []string
	err = json.Unmarshal(projectVersionsJson, &versions)
	if err != nil {
		logrus.Errorln(err)
	}

	for _, version := range versions {
		// get release notes of specific release
		pluginJson, _ := r.DB.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, version))
		var releaseNote types.ReleaseNote
		err := json.Unmarshal(pluginJson, &releaseNote)
		if err != nil {
			logrus.Errorln(err)
			// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
			return releaseNotes, errors.New("failed to unmarshal ReleaseNote")
		}
		// append release version note to list
		// releaseNote.Body = releaseNote.BodyHTML
		releaseNotes = append(releaseNotes, releaseNote)

	}

	return releaseNotes, nil
}
