package pluginManager

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/manifest"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/updateCenter/pluginVersions"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/sirupsen/logrus"
)

type PluginType int

const (
	UNKNOWN PluginType = iota
	EXTERNAL
	TRANSITIVE // used as dependency for some plugins
	PRIMARY    // no one rely on it as dep
)

func (pt PluginType) String() string {
	return []string{"Unknown", "External", "Transitive", "Primary"}[pt]
}

// func NewDependency(name string, version string, optional bool) Dependency {
// 	return Dependency{Name: name, Version: version, Optional: optional}
// }

// https://github.com/emelianrus/jenkins-update-center/blob/master/pkg/updateCenter/updateCenter.go#L97
type Warnings struct {
	Id      string
	Message string
	Name    string
	// Type     string
	Url      string
	Versions []struct {
		LastVersion string
		Pattern     string
	}
}

type Plugin struct {
	Name    string
	Version string

	Url                 string
	Type                PluginType
	RequiredCoreVersion string
	Dependencies        map[string]Plugin
	RequiredBy          map[string]string
	Optional            bool // rely on parent plugin
	Warnings            []Warnings
}

// Create Plugin from name and version
func NewPluginWithVersion(name string, version string) *Plugin {
	logrus.Debugf("Creating new plugin: '%s' with version: '%s'\n", name, version)

	pl := &Plugin{
		Name:         name,
		Version:      version,
		Url:          "default-value-for-url",
		Type:         UNKNOWN,
		Dependencies: make(map[string]Plugin),
		RequiredBy:   make(map[string]string),
	}
	return pl
}

// Create Plugin from name and url
func NewPluginWithUrl(name string, url string) *Plugin {
	logrus.Debugf("Creating new plugin: '%s' with url: '%s'\n", name, url)
	pl := &Plugin{
		Name:         name,
		Version:      "",
		Url:          url,
		Type:         UNKNOWN,
		Dependencies: make(map[string]Plugin),
		RequiredBy:   make(map[string]string),
	}
	return pl
}

const JENKINS_PLUGINS_URL = "https://updates.jenkins.io"

// download plugin hpi file from jenkins update center
func (p *Plugin) Download() (string, error) {

	/*
		if we dont have version but have url
			* we need to download plugin and read manifest from hpi file
			manifest, _ := manifest.Parse(filename)
		if we have version but dont have url
			* we need to get url during download
	*/
	// External plugin has URL to download but doesn't have version field

	logrus.Infof("Downloading plugin with params:: name: %s, ver: %s, url: %s\n", p.Name, p.Version, p.Url)

	var isExternalPluginType bool = false

	if p.Version == "" {
		isExternalPluginType = true
	}

	var fileLocation string

	if !isExternalPluginType {
		p.Url = fmt.Sprintf("%s/download/plugins/%s/%s/%s.hpi", JENKINS_PLUGINS_URL, p.Name, p.Version, p.Name)
		// where to store downloaded file
		fileLocation = fmt.Sprintf("plugins/%s-%s.hpi", p.Name, p.Version)
	}

	if !utils.IsFileExist("plugins") {
		// Create plugins dir where to store cache plugins
		err := os.Mkdir("plugins", os.ModePerm)
		if err != nil {
			logrus.Error(err)
			return "", errors.New("can't create dir")
		}
	}

	if utils.IsFileExist(fileLocation) {
		logrus.Debugf("File %s Exist. skipped download\n", fileLocation)
	} else {
		logrus.Infof("Downloading plugin %s\n", p.Name)
		response := utils.DoRequestGet(p.Url)
		// Create the file
		file, err := os.Create(fileLocation)
		if err != nil {
			logrus.Error(err)
			return "", errors.New("can't create file")
		}
		defer file.Close()

		// write content to file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			logrus.Error(err)
			return "", errors.New("can't write content to file")
		}

		defer response.Body.Close()
		logrus.Infof("File downloaded to %s", fileLocation)
	}

	if isExternalPluginType {
		manifestFile, _ := manifest.Parse(fileLocation)
		p.Version = manifestFile["Plugin-Version"]
	}

	logrus.Infof("Downloaded Plugin name: %s, Plugin version: %s, Plugin URL: %s", p.Name, p.Version, p.Url)

	return fileLocation, nil
}

// Loads dependencies from jenkins update center into Plugin struct
func (p *Plugin) LoadDependenciesFromUpdateCenter() map[string]Plugin {
	pv, _ := pluginVersions.Get()
	for _, dep := range pv.Plugins[p.Name][p.Version].Dependencies {
		if !dep.Optional {

			p.Dependencies[dep.Name] = Plugin{
				Name:    dep.Name,
				Version: dep.Version,
			}

		}
	}
	return p.Dependencies
}

// Loads dependencies from hpi file manifest into Plugin struct
func (p *Plugin) LoadDependenciesFromManifest() map[string]Plugin {
	p.Download()
	// we need manifest to get jenkins core version to get the right update center json
	manifestFile, _ := manifest.Parse(fmt.Sprintf("plugins/%s-%s.hpi", p.Name, p.Version))

	logrus.Debugf("[GetDependenciesFromManifest] plugin name: %s jenkins core: %s\n\n", p.Name, manifestFile["Jenkins-Version"])

	for _, dep := range manifestFile.GetDependencies() {
		logrus.Debugf("[GetDependenciesFromManifest] %s: all deps from manifest Name: %s Version: %s Optional: %t", p.Name, dep.Name, dep.Version, dep.Optional)
		if !dep.Optional {
			p.Dependencies[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)
		}
	}
	logrus.Debugf("[GetDependenciesFromManifest] Plugin: %s has deps: %v", p.Name, p.Dependencies)
	return p.Dependencies
}

func (p *Plugin) GetManifestAttrs() map[string]string {
	p.Download()
	// we need manifest to get jenkins core version to get the right update center json
	manifestFile, _ := manifest.Parse(fmt.Sprintf("plugins/%s-%s.hpi", p.Name, p.Version))

	attrs := make(map[string]string)
	for k, v := range manifestFile {
		attrs[k] = v
	}

	return attrs
}
