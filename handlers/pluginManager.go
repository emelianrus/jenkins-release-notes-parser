package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/pluginManager"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TODO: need to know what we need to rescan :)
func (s *ProjectService) RescanProjectNow(c *gin.Context) {
	logrus.Infoln("RescanProjectNow route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Received request body: %+v\n", body)

	c.String(http.StatusOK, "RescanProjectNow")
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

	c.String(http.StatusOK, "AddPluginsFile")
}

// Plugin-manager handler to add new plugin to plugin-manager
func (s *ProjectService) AddNewPlugin(c *gin.Context) {
	logrus.Infoln("AddNewPlugin route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Received request body: %+v\n", body)

	// check if plugin is not exist you can not add it to list
	if _, exists := s.PluginManager.PluginVersions.Plugins[body["name"]]; !exists {
		c.String(http.StatusBadRequest, fmt.Sprintf("AddNewPlugin %s:%s is not exist in plugins", body["name"], body["version"]))
		return
	}

	s.PluginManager.AddPluginWithVersion(body["name"], body["version"])

	c.String(http.StatusOK, fmt.Sprintf("AddNewPlugin %s:%s", body["name"], body["version"]))
}

func (s *ProjectService) DeletePlugin(c *gin.Context) {
	logrus.Infoln("DeletePlugin route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Received request body: %+v\n", body)

	s.PluginManager.DeletePlugin(body["name"])

	c.String(http.StatusOK, fmt.Sprintf("DeletePlugin %s", body["name"]))
}

func (s *ProjectService) EditCoreVersion(c *gin.Context) {
	logrus.Infoln("[EditCoreVersion] route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("[EditCoreVersion] Received request body: %+v\n", body)

	s.PluginManager.SetCoreVersion(body["name"])

	logrus.Infof("[EditCoreVersion] set new coreVersion: %s\n", body["name"])
	c.String(http.StatusOK, fmt.Sprintf("EditCoreVersion %s", body["name"]))
}

func (s *ProjectService) GetCoreVersion(c *gin.Context) {
	logrus.Infoln("GetCoreVersion route reached")
	c.String(http.StatusOK, s.PluginManager.GetCoreVersion())
}

func (s *ProjectService) GetPotentialUpdates(c *gin.Context) {
	logrus.Infoln("GetPotentialUpdates route reached")
	potentialUpdates, _ := s.Redis.GetPotentialUpdates()
	c.JSON(http.StatusOK, potentialUpdates)
}

func (s *ProjectService) CheckDeps(c *gin.Context) {
	logrus.Infoln("CheckDeps route reached")

	c.JSON(http.StatusOK, s.PluginManager.FixPluginDependencies())
}
func (s *ProjectService) GetPluginsData(c *gin.Context) {
	logrus.Infoln("GetPluginsData route reached")

	type pluginManagerData struct {
		Plugins     map[string]*pluginManager.Plugin
		CoreVersion string
	}

	data := pluginManagerData{
		Plugins:     s.PluginManager.GetPlugins(),
		CoreVersion: s.PluginManager.GetCoreVersion(),
	}
	c.JSON(http.StatusOK, data)

}
func (s *ProjectService) GetFixedDepsDiff(c *gin.Context) {
	logrus.Infoln("GetFixedDepsDiff route reached")

	c.JSON(http.StatusOK, s.PluginManager.GetFixedDepsDiff())
}
