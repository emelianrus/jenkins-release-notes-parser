package handlers

import (
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/pluginManager"
	"github.com/emelianrus/jenkins-release-notes-parser/storage/redisStorage"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// struct for handlers to use DB connection
type ProjectService struct {
	Redis         *redisStorage.RedisStorage
	PluginManager pluginManager.PluginManager
}

func (s *ProjectService) GetAllProjects(c *gin.Context) {
	logrus.Infoln("GetAllProjects route reached")
	projects, _ := s.Redis.GetAllProjects()
	watcherList, _ := s.Redis.GetWatcherData()

	type resultPlugins struct {
		IsInWatcherList bool
		Project         types.Project
	}
	result := []resultPlugins{}

	for _, project := range projects {
		var inWatcherList bool = false
		if _, ok := watcherList[project.Name]; ok {
			inWatcherList = true
		}

		result = append(result, resultPlugins{
			IsInWatcherList: inWatcherList,
			Project:         project,
		})
	}

	c.JSON(http.StatusOK, result)
}

// TODO: https://api.github.com/repos/OWNER/REPO/releases
func (s *ProjectService) GetProjectReleaseNotes(c *gin.Context) {
	logrus.Infoln("GetProjectReleaseNotes route reached")
	ownerName := c.Param("owner")
	repoName := c.Param("repo")

	releaseNotes, err := s.Redis.GetProjectReleaseNotes(ownerName, repoName)
	if err != nil {
		logrus.Errorf("can not get project %s:%s\n", ownerName, repoName)
		logrus.Errorln(err)
	}

	type projectReleaseNotes struct {
		Repo         string
		Owner        string
		ProjectGroup string
		ReleaseNotes []types.ReleaseNote
	}

	resp := projectReleaseNotes{
		Repo:         repoName,
		Owner:        ownerName,
		ProjectGroup: "jenkinsci",
		ReleaseNotes: releaseNotes,
	}
	c.JSON(http.StatusOK, resp)
	logrus.Infof("HITED %s/%s\n", ownerName, repoName)
}
