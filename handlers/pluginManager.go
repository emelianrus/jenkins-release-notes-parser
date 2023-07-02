package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *ProjectService) RescanProjectNow(c *gin.Context) {
	logrus.Infoln("RescanProjectNow route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Received request body: %+v\n", body)

	res := s.PluginManager.GetPlugins() //[body["name"]].Download()

	res[body["name"]].Download()

	c.String(http.StatusOK, "RescanProjectNow")
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
