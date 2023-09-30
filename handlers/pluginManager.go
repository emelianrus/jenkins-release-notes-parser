package handlers

import (
	"fmt"
	"net/http"
	"os"

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
	if _, pluginExists := s.PluginManager.PluginVersions.Plugins[body["name"]]; !pluginExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("AddNewPlugin %s:%s is not exist in public plugins", body["name"], body["version"]),
			"status":  "error",
		})
		return
	}

	if _, ok := s.PluginManager.PluginVersions.Plugins[body["name"]][body["version"]]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("AddNewPlugin %s:%s exist in public plugins but it doesnt have version specified", body["name"], body["version"]),
			"status":  "error",
		})
		return
	}

	s.PluginManager.AddPluginWithVersion(body["name"], body["version"])

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("AddNewPlugin %s:%s", body["name"], body["version"]),
		"status":  "ok",
	})
}

func (s *ProjectService) GetManifestAttrs(c *gin.Context) {
	logrus.Infoln("GetManifestAttrs route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Received request body: %+v\n", body)

	attrs, err := s.PluginManager.GetPlugin(body["name"]).GetManifestAttrs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attrs)
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

func (s *ProjectService) CheckDeps(c *gin.Context) {
	logrus.Infoln("CheckDeps route reached")
	c.JSON(http.StatusOK, s.PluginManager.FixPluginDependencies())
}
func (s *ProjectService) GetPluginsData(c *gin.Context) {
	logrus.Infoln("GetPluginsData route reached")

	type pluginManagerData struct {
		Plugins             map[string]*pluginManager.Plugin
		CoreVersion         string
		UpdateCenterVersion string
	}

	c.JSON(http.StatusOK, pluginManagerData{
		Plugins:             s.PluginManager.GetPlugins(),
		CoreVersion:         s.PluginManager.GetCoreVersion(),
		UpdateCenterVersion: s.PluginManager.GetUpdateCenterVersion(),
	})
}
func (s *ProjectService) GetReleaseNotesDiff(c *gin.Context) {
	logrus.Infoln("GetReleaseNotesDiff route reached")

	type ProjectData struct {
		Name     string   `json:"name"`
		Versions []string `json:"versions"`
	}
	var projectData ProjectData
	if err := c.BindJSON(&projectData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use the existing context in your function
	result := make(chan interface{})
	go func() {
		defer close(result)
		result <- s.PluginManager.GetReleaseNotesDiff(projectData.Name, projectData.Versions)
	}()

	// Listen for the client connection close signal
	select {
	case <-c.Writer.CloseNotify():
		logrus.Infoln("Client connection closed")
		// Perform any cleanup or cancellation logic here
		return
	case data := <-result:
		c.JSON(http.StatusOK, data)
	}
}

func (s *ProjectService) GetVersionsDiff(c *gin.Context) {
	logrus.Infoln("GetVersionsDiff route reached")

	c.JSON(http.StatusOK, s.PluginManager.GetFixedDepsDiff())
}

func (s *ProjectService) DownloadFilePluginManager(c *gin.Context) {
	// TODO: payload txt file or yaml file or any other type
	logrus.Infoln("DownloadFile route reached")

	data := s.PluginManager.GenerateFileOutputPluginManager()
	if len(data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "doesn't have data to create file for plugin manager",
			"status":  "error",
		})
		return
	}
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "file*.txt")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create temp file")
		return
	}
	defer tmpFile.Close()

	// Write the contents to the file
	_, err = tmpFile.Write(data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to write to temp file")
		return
	}

	// Set the appropriate headers
	c.Header("Content-Disposition", "attachment; filename=file.txt")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	// Serve the file
	c.File(tmpFile.Name())
}

// TODO: DRY
func (s *ProjectService) DownloadFilePluginChanges(c *gin.Context) {
	// TODO: payload txt file or yaml file or any other type
	logrus.Infoln("DownloadFile route reached")
	data := s.PluginManager.GenerateFileOutputUpdatedPlugins()
	if len(data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "doesn't have data to create file for plugin changes",
			"status":  "error",
		})
		return
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "file*.txt")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create temp file")
		return
	}
	defer tmpFile.Close()

	// Write the contents to the file
	_, err = tmpFile.Write(data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to write to temp file")
		return
	}

	// Set the appropriate headers
	c.Header("Content-Disposition", "attachment; filename=file.txt")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	// Serve the file
	c.File(tmpFile.Name())
}
