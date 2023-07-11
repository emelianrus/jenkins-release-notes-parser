package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *ProjectService) GetUpdatedPluginList(c *gin.Context) {
	logrus.Infoln("GetUpdatedPluginList route reached")
	updatedPluginList := make(map[string]string)

	for _, pl := range s.PluginManager.UpdatedPlugins {
		updatedPluginList[pl.Name] = pl.Version
	}
	c.JSON(http.StatusOK, updatedPluginList)
}

func (s *ProjectService) AddUpdatedPluginList(c *gin.Context) {
	logrus.Infoln("AddUpdatedPluginList route reached")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plugins := s.PluginManager.FileParser.Parse(body)

	// clean what we have
	for k := range s.PluginManager.UpdatedPlugins {
		delete(s.PluginManager.UpdatedPlugins, k)
	}

	for pluginName, version := range plugins {
		s.PluginManager.SetUpdatedPluginWithVersion(pluginName, version)
	}

	c.String(http.StatusOK, "AddUpdatedPluginList")
}
