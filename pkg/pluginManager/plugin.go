package pluginManager

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/manifest"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/updateCenter/pluginVersions"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/updateCenter/updateCenter"
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
	return &Plugin{
		Name:         name,
		Version:      version,
		Url:          "",
		Type:         UNKNOWN,
		Dependencies: make(map[string]Plugin),
		RequiredBy:   make(map[string]string),
	}
}

// Create Plugin from name and url
func NewPluginWithUrl(name string, url string) *Plugin {
	logrus.Debugf("Creating new plugin: '%s' with url: '%s'\n", name, url)
	return &Plugin{
		Name:         name,
		Version:      "",
		Url:          url,
		Type:         UNKNOWN,
		Dependencies: make(map[string]Plugin),
		RequiredBy:   make(map[string]string),
	}
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

// Loads warning to plugin struct
func (p *Plugin) LoadWarnings() {
	// download plugin hpi file
	p.Download()
	// we need manifest to get jenkins core version to get the right update center json
	manifestFile, _ := manifest.Parse(fmt.Sprintf("plugins/%s-%s.hpi", p.Name, p.Version))
	// get update center for current plugin, we will get warnings from UC
	uc, _ := updateCenter.Get(manifestFile["Jenkins-Version"])

	// clear warnings but keep allocated memory
	p.Warnings = p.Warnings[:0]

	for _, warn := range uc.Warnings {

		// skip all plugins except current one
		if warn.Name == p.Name {

			for _, warningVersion := range warn.Versions {
				// if warning version lower then current than we dont have warning
				if warningVersion.LastVersion == p.Version || utils.IsNewerThan(warningVersion.LastVersion, p.Version) {
					// write errors back to plugin
					// TODO: rewrite resolve deps function, may reuse warnings list when newer plugin comes
					p.Warnings = append(p.Warnings, Warnings{
						Id:      warn.Id,
						Message: warn.Message,
						Name:    warn.Name,
						Url:     warn.Url,
						Versions: []struct {
							LastVersion string
							Pattern     string
						}(warn.Versions),
					})
				}
			}
		}
	}
}

// get current plugin required core version
func (p *Plugin) LoadRequiredCoreVersion() {
	pv, _ := pluginVersions.Get()
	fmt.Println(pv.Plugins[p.Name][p.Version].RequiredCore)
	requiredCore := pv.Plugins[p.Name][p.Version].RequiredCore
	if requiredCore == "" {
		logrus.Warnln("[LoadRequiredCoreVersion] version is empty")
	}
	p.RequiredCoreVersion = pv.Plugins[p.Name][p.Version].RequiredCore

}

// should fix warnings and paste back to struct
// TODO: what should we do if all versions with error???
func (p *Plugin) FixWarnings() error {
	// TODO: do not call each time reload warnings everywhere
	p.LoadWarnings() // we need to find last version with warning
	if len(p.Warnings) == 0 {
		logrus.Debugf("No error found for plugin %s version %s", p.Name, p.Version)
		return nil
	}
	pv, _ := pluginVersions.Get() // to check if there any version where warning version +1
	fmt.Println(pv.Plugins[p.Name][p.Version].RequiredCore)

	var versions []string

	for version := range pv.Plugins[p.Name] {
		versions = append(versions, version)
	}

	// TODO: refact
	var nextVersion string
	var currentVersion string
	currentVersion = p.Version
	var nPl Plugin
	for {
		var err error
		nextVersion, err = utils.GetNextVersion(versions, currentVersion)

		if nextVersion == "" {
			logrus.Infoln("Reach get version limit")
		}

		if err != nil {
			logrus.Warnf("Error during GetNextVersion for %s\n", currentVersion)
			logrus.Warnln(err)

			return err
		}
		currentVersion = nextVersion

		nPl = Plugin{
			Name:    p.Name,
			Url:     p.Url,
			Version: nextVersion,
		}

		nPl.LoadWarnings()

		fmt.Printf("pl : %s version %s warn: %d", nPl.Name, nPl.Version, len(nPl.Warnings))

		if len(nPl.Warnings) == 0 {
			break
		}
	}

	// set new version as current and remove deps as could differs
	p.Version = nPl.Version
	p.Url = nPl.Url
	p.Warnings = nPl.Warnings
	p.Dependencies = make(map[string]Plugin)
	// TODO: error if latest version also has error
	return nil
}

// Loads dependencies from jenkins update center into Plugin struct
func (p *Plugin) LoadDependenciesFromUpdateCenter() map[string]Plugin {
	pv, _ := pluginVersions.Get()
	for _, dep := range pv.Plugins[p.Name][p.Version].Dependencies {
		if !dep.Optional {
			p.Dependencies[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)
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

// TODO: fix me
// // returns map [plugin name] = warning
// func (p *Plugin) CheckWarnings(coreVersion string) map[string][]string {
// 	plugins := make(map[string]Plugin)

// 	for _, plugin := range p.Plugins {
// 		plugins[plugin.Name] = *plugin
// 	}

// 	uc := updateCenter.Get(coreVersion)
// 	resultWarnings := make(map[string][]string)

// 	for _, warn := range uc.Warnings {
// 		for _, plugin := range plugins {
// 			if warn.Name == plugin.Name {
// 				for _, pluginwarn := range warn.Versions {
// 					var warnVersion = regexp.MustCompile(pluginwarn.Pattern)
// 					if warnVersion.MatchString(plugin.Version) {
// 						message := fmt.Sprintf("%s\n  LastVersions: %s\n  Pattern: %s\n  URL %s\n", warn.Message, pluginwarn.LastVersion, pluginwarn.Pattern, warn.Url)
// 						resultWarnings[plugin.Name+":"+plugin.Version] = append(resultWarnings[plugin.Name+plugin.Version], message)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return resultWarnings
// }

// func (p *Plugin) getDependenciesFromMainfest() []manifest.Dependency {

// 	filename := fmt.Sprintf("plugins/%s-%s.hpi", p.Name, p.Version)
// 	p.Download()
// 	manifest, _ := manifest.Parse(filename)
// 	logrus.Infof("plugin name: %s:%s\n", p.Name, p.Version)

// var wg sync.WaitGroup
// // download deps in parallel
// for _, dep := range deps {
// 	wg.Add(1)
// 	dn := dep.Name
// 	dv := dep.Version
// 	go func() {
// 		defer wg.Done()
// 		plugin := Plugin{Name: dn, Version: dv}
// 		plugin.Download()
// 	}()
// }
// wg.Wait()
// 	return manifest.GetDependencies()
// }
