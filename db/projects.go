package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
	"github.com/sirupsen/logrus"
)

func (r *Redis) GetWatcherProjects() ([]types.Project, error) {
	var projects []types.Project
	// TODO: this should be configurable
	repoOwner := "jenkinsci"

	watcherProjects, _ := r.GetWatcherData()

	// gather project data
	for projectName := range watcherProjects {

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

func (r *Redis) GetAllProjects() ([]types.Project, error) {
	var projects []types.Project
	// TODO: this should be configurable
	repoOwner := "jenkinsci"

	// get list of all projects
	var projectsTmp []string

	projectsKeys, _ := r.Keys(fmt.Sprintf("github:%s:*", repoOwner))

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
func (r *Redis) GetProjectReleaseNotes(projectOwner string, projectName string) ([]types.ReleaseNote, error) {

	releaseNotes := []types.ReleaseNote{}

	// get all versions for specific project
	projectVersionsJson, err := r.Get(fmt.Sprintf("github:%s:%s:versions", projectOwner, projectName)).Bytes()
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
		pluginJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, version)).Bytes()
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

type PotentialUpdates map[string][]types.ReleaseNote

func (r *Redis) GetPotentialUpdates() (PotentialUpdates, error) {
	watcherListProjects, _ := r.GetWatcherData()
	potentialUpdates := PotentialUpdates{}
	cachedProjects, _ := r.GetAllProjects()

	for watcherProjectName, watcherProjectVersion := range watcherListProjects {

		for _, cachedProject := range cachedProjects {
			if cachedProject.Name == watcherProjectName {
				// we hit cached plugin with watcher plugin name
				// now need to get release notes from top of cached to watcher set version

				releaseNotes, _ := r.GetProjectReleaseNotes("jenkinsci", watcherProjectName)
				// iterate over release notes

				var resultRelaseNotes []types.ReleaseNote

				for _, releaseNote := range releaseNotes {
					// reached installed version. break
					logrus.Debugf("releaseNote.Name %s ||| watcherProject.Version %s\n", releaseNote.Name, watcherProjectVersion)
					if releaseNote.Name == watcherProjectVersion {
						break
					}

					// if not reached
					resultRelaseNotes = append(resultRelaseNotes, releaseNote)
				}
				if len(resultRelaseNotes) > 0 {
					potentialUpdates[watcherProjectName] = resultRelaseNotes
				}

			}

		}

	}
	fmt.Println(potentialUpdates)
	return potentialUpdates, nil

}
