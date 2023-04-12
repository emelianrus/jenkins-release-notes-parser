package handlers

import (
	"fmt"
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// struct for handlers to use DB connection
type ProjectService struct {
	Redis *db.Redis
}

func (s *ProjectService) AddProject(c *gin.Context) {
	logrus.Infoln("AddProject")
	c.JSON(http.StatusOK, "ASd")
}

func (s *ProjectService) AddMultiplyProjects(c *gin.Context) {
	logrus.Infoln("AddMultiplyProjects route reached")
	c.JSON(http.StatusOK, "ASd")
}

func (s *ProjectService) GetProjectById(c *gin.Context) {
	logrus.Infoln("GetProjectById route reached")
	projectName := c.DefaultQuery("name", "")
	c.String(http.StatusOK, "Hello %s", projectName)
}

func (s *ProjectService) GetAllProjects(c *gin.Context) {
	logrus.Infoln("GetAllProjects route reached")
	projects, _ := s.Redis.GetAllProjects()
	c.JSON(http.StatusOK, projects)
}

func (s *ProjectService) GetProjectsById(c *gin.Context) {
	logrus.Infoln("GetProjectsById route reached")
	c.JSON(http.StatusOK, "GetProjectsById")
}

func (s *ProjectService) DeleteProject(c *gin.Context) {
	logrus.Infoln("DeleteProject route reached")
	id := c.DefaultQuery("id", "")
	c.String(http.StatusOK, "Hello %s", id)
}

func (s *ProjectService) GetWatcherList(c *gin.Context) {
	watcherList, err := s.Redis.GetWatcherList()
	if err != nil {
		logrus.Errorln("can not get watcher list")
		logrus.Errorln(err)
	}
	c.JSON(http.StatusOK, watcherList)
}
func (s *ProjectService) EditWatcherList(c *gin.Context) {
	logrus.Infoln("EditWatcherList route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Received request body: %+v\n", body)

	err := s.Redis.SetWatcherList(body)
	if err != nil {
		logrus.Errorln("can not set watcher list to DB")
		logrus.Errorln(err)
	}

	c.String(http.StatusOK, "EditWatcherList")
}

func (s *ProjectService) DeleteMultiplyProjects(c *gin.Context) {
	logrus.Infoln("DeleteMultiplyProjects route reached")
	var ids []string
	if err := c.BindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	// Your logic to delete the items with the given IDs goes here
	c.JSON(http.StatusOK, gin.H{"message": "Deleted items with IDs", "ids": ids})
}

// TODO: https://api.github.com/repos/OWNER/REPO/releases
func (s *ProjectService) GetProjectReleaseNotes(c *gin.Context) {
	logrus.Infoln("GetProjectReleaseNotes route reached")
	ownerName := c.Param("owner")
	repoName := c.Param("repo")

	releaseNotes, err := s.Redis.GetProject(ownerName, repoName)
	if err != nil {
		logrus.Errorf("can not get project %s:%s\n", ownerName, repoName)
		logrus.Errorln(err)
	}

	type Resp struct {
		Repo         string
		Owner        string
		ProjectGroup string
		ReleaseNotes []types.ReleaseNote
	}

	resp := Resp{
		Repo:         repoName,
		Owner:        ownerName,
		ProjectGroup: "jenkinsci",
		ReleaseNotes: releaseNotes,
	}
	c.JSON(http.StatusOK, resp)
	fmt.Printf("HITED %s/%s\n", ownerName, repoName)
}
