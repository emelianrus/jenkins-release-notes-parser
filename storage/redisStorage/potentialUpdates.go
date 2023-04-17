package redisStorage

import (
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/sirupsen/logrus"
)

type PotentialUpdates map[string][]types.ReleaseNote

func (r *RedisStorage) GetPotentialUpdates() (PotentialUpdates, error) {
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

				foundVersion := false
				for _, releaseNote := range releaseNotes {

					// reached installed version. break
					// logrus.Debugf("releaseNote.Name %s ||| watcherProject.Version %s\n", releaseNote.Name, watcherProjectVersion)
					if releaseNote.Name == watcherProjectVersion {
						foundVersion = true
						break
					}

					// if not reached
					resultRelaseNotes = append(resultRelaseNotes, releaseNote)
				}
				if !foundVersion {
					logrus.Warnf("haven't found version %s:%s\n", cachedProject.Name, watcherProjectVersion)
					// TODO: try to download from github
					// add feed
				}

				if len(resultRelaseNotes) > 0 {
					potentialUpdates[watcherProjectName] = resultRelaseNotes
				}

			}

		}

	}

	return potentialUpdates, nil

}
