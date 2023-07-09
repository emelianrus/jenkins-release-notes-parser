package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *ProjectService) GetPluginList(c *gin.Context) {
	watcherList, err := s.Redis.GetPluginListData()
	if err != nil {
		logrus.Errorln("can not get watcher list")
		logrus.Errorln(err)
	}
	c.JSON(http.StatusOK, watcherList)
}

func (s *ProjectService) AddPluginsFile(c *gin.Context) {
	logrus.Infoln("AddPluginsFile route reached")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plugins := s.PluginManager.FileParser.Parse(body)

	for pluginName, version := range plugins {
		s.PluginManager.AddPluginWithVersion(pluginName, version)
	}

	s.Redis.SetPluginListData(plugins)

	c.String(http.StatusOK, "AddPluginsFile")
}

func (s *ProjectService) EditWatcherList(c *gin.Context) {
	logrus.Infoln("EditWatcherList route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Received request body: %+v\n", body)

	err := s.Redis.SetPluginListData(body)
	if err != nil {
		logrus.Errorln("can not set watcher list to DB")
		logrus.Errorln(err)
	}

	c.String(http.StatusOK, "EditWatcherList")
}
