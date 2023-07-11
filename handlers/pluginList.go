package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *ProjectService) GetPluginList(c *gin.Context) {
	resp := make(map[string]string)
	for _, pl := range s.PluginManager.Plugins {
		resp[pl.Name] = pl.Version
	}
	c.JSON(http.StatusOK, resp)
}

func (s *ProjectService) AddPluginsFile(c *gin.Context) {
	logrus.Infoln("AddPluginsFile route reached")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plugins := s.PluginManager.FileParser.Parse(body)

	s.PluginManager.CleanPlugins()

	for pluginName, version := range plugins {
		s.PluginManager.AddPluginWithVersion(pluginName, version)
	}

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

	s.PluginManager.CleanPlugins()

	for k, v := range body {
		s.PluginManager.AddPluginWithVersion(k, v)
	}

	err := s.Redis.SetPluginListData(body)
	if err != nil {
		logrus.Errorln("can not set watcher list to DB")
		logrus.Errorln(err)
	}

	c.String(http.StatusOK, "EditWatcherList")
}
