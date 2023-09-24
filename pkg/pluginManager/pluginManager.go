package pluginManager

import (
	"fmt"
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/outputGenerators"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/parsers"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/sources"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/sources/github"
	jenkins "github.com/emelianrus/jenkins-release-notes-parser/pkg/sources/jenkinsPluginSite"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-update-center/pkg/jenkinsSite"
	"github.com/sirupsen/logrus"
)

type PluginManager struct {
	Plugins map[string]*Plugin // Main list where we have plugins. We use map to speedup get set find

	UpdatedPlugins map[string]*Plugin // Temp storage. will set this after update/fix deps

	coreVersion string // jenkins core version

	JenkinsSite    jenkinsSite.JenkinsSiteInterface // external. jenkins site main struct
	UpdateCenter   jenkinsSite.UpdateCenter         // external. from jenkins api. plugin latest version + deprecations + deps for plugin from JenkinsSite
	PluginVersions jenkinsSite.PluginVersions       // external. from jenkins api. plugins data with versions from JenkinsSite

	// release notes sources
	PluginSite   jenkins.PluginSite // jenkins site which has release notes (last 10)
	GitHubClient github.GitHub      // github client to get release notes

	// file read/write logic
	FileParser parsers.InputParser            // parses txt file into plugin-manager plugins
	FileOutput outputGenerators.FileGenerator // plugin-manager plugins into file content
}

func NewPluginManager() PluginManager {

	js := jenkinsSite.NewJenkinsSite()
	stableCoreVersion, _ := js.GetStableCoreVersion()
	uc, _ := js.GetUpdateCenter(stableCoreVersion)
	pv, _ := js.GetPluginVersions()

	return PluginManager{
		coreVersion: stableCoreVersion,

		Plugins:        make(map[string]*Plugin),
		UpdatedPlugins: make(map[string]*Plugin),

		JenkinsSite:    js,
		UpdateCenter:   *uc,
		PluginVersions: *pv,

		PluginSite:   jenkins.NewPluginSite(),
		GitHubClient: github.NewGitHubClient(),

		// in/out from plugin manager to file
		FileParser: parsers.TXTFileParser{},
		FileOutput: outputGenerators.TxtOutput{},
	}
}

func (pm *PluginManager) SetFileParser(parser parsers.InputParser) {
	pm.FileParser = parser
}

func (pm *PluginManager) SetJenkinsSite(js jenkinsSite.JenkinsSiteInterface) {
	pm.JenkinsSite = js
}
func (pm *PluginManager) SetUpdateCenter(uc *jenkinsSite.UpdateCenter) {
	pm.UpdateCenter = *uc
}
func (pm *PluginManager) SetPluginVersions(pv *jenkinsSite.PluginVersions) {
	pm.PluginVersions = *pv
}

func (pm *PluginManager) SetFileOutput(og outputGenerators.FileGenerator) {
	pm.FileOutput = outputGenerators.SetOutputGenerator(og)
}

func (pm *PluginManager) GenerateFileOutputUpdatedPlugins() []byte {
	result := make(map[string]string)

	for _, pl := range pm.UpdatedPlugins {
		result[pl.Name] = pl.Version
	}
	return pm.FileOutput.Generate(result)
}
func (pm *PluginManager) GenerateFileOutputPluginManager() []byte {
	result := make(map[string]string)

	for _, pl := range pm.Plugins {
		result[pl.Name] = pl.Version
	}
	return pm.FileOutput.Generate(result)
}

// iterates over all plugin and add requiredBy
func (pm *PluginManager) generateRequiredBy() {
	logrus.Infoln("[generateRequiredBy]")

	for _, pl := range pm.Plugins {
		for depName := range pl.Dependencies {
			if _, exists := pm.Plugins[depName]; exists {
				pm.Plugins[depName].RequiredBy[pl.Name] = pl.Version
			}
			// else {
			// 	delete(pm.Plugins[depName].RequiredBy, pl.Name)
			// }
		}
	}
}

func (pm *PluginManager) cleanRequiredBy(removedProjectName string) {
	logrus.Infoln("[cleanRequiredBy]")

	for _, pl := range pm.Plugins {
		if _, exists := pl.RequiredBy[removedProjectName]; exists {
			logrus.Infof("removed plugin %s from %s.RequiredBy\n", pl.Name, removedProjectName)
			delete(pl.RequiredBy, removedProjectName)
		}
	}
}

// load latest version/url/reuiredcore/deps
func (pm *PluginManager) reloadPluginsData(p *Plugin) {
	logrus.Infof("reloadPluginsData %s:%s", p.Name, p.Version)
	// NOTE: if version is incorrect will not get any data
	for _, dep := range pm.PluginVersions.Plugins[p.Name][p.Version].Dependencies {
		if !dep.Optional {
			p.Dependencies[dep.Name] = Plugin{
				Name:    dep.Name,
				Version: dep.Version,
			}
		}
	}
	// Load url
	p.Url = pm.PluginVersions.Plugins[p.Name][p.Version].Url

	pluginName := p.Name
	if !strings.HasSuffix(pluginName, "-plugin") {
		pluginName = pluginName + "-plugin"
	}
	p.GITUrl = "https://github.com/jenkinsci/" + pluginName

	// Load required core version
	p.RequiredCoreVersion = pm.PluginVersions.Plugins[p.Name][p.Version].RequiredCore

	// get plugin latest version, doesnt have endpoint with latest version so need to detect
	latestVersion := "0.0.0" // set lowest version possible
	for version := range pm.PluginVersions.Plugins[p.Name] {
		if utils.IsNewerThan(version, latestVersion) {
			latestVersion = version
		}
	}

	p.LatestVersion = latestVersion

}

func (pm *PluginManager) GetPlugins() map[string]*Plugin {
	return pm.Plugins
}
func (pm *PluginManager) GetPlugin(name string) *Plugin {
	return pm.Plugins[name]
}

func (pm *PluginManager) GetCoreVersion() string {
	return pm.coreVersion
}
func (pm *PluginManager) SetCoreVersion(newCoreVersion string) {
	pm.coreVersion = newCoreVersion
	uc, _ := pm.JenkinsSite.GetUpdateCenter(newCoreVersion)

	pm.UpdateCenter = *uc
}

func (pm *PluginManager) AddPluginWithUrl(pluginName string, url string) {
	// TODO
}

func (pm *PluginManager) AddPlugin(plugin *Plugin) {
	pm.Plugins[plugin.Name] = plugin
}

func (pm *PluginManager) LoadPluginData(p *Plugin) {
	pm.reloadPluginsData(p)
}

func (pm *PluginManager) AddPluginWithVersion(pluginName string, version string) {
	logrus.Infof("[AddPluginWithVersion] Adding new plugin to pluginManager %s:%s", pluginName, version)
	pl := NewPluginWithVersion(pluginName, version)

	logrus.Infof("Adding new plugin to pluginManager %s:%s", pl.Name, pl.Version)
	if _, found := pm.Plugins[pl.Name]; found {
		logrus.Warnf("Found copy of plugin in pluginsfile. Plugin name: '%s'\n", pl.Name)
	}

	pm.reloadPluginsData(pl)
	// TODO: sync with DB

	pm.AddPlugin(pl)

	// functions should be executed after each plugin added
	// requires plugin exist in pm.Plugins[pl.Name]
	pm.generateRequiredBy()
	pm.LoadWarnings()
}

// will try to delete plugin from list if no one rely as dep on this plugin
func (p *PluginManager) DeletePlugin(pluginName string) {
	// check if we can delete plugin and no one rely on it
	// if some one use it print
	logrus.Infof("[DeletePlugin] %s\n", pluginName)
	// clean removed plugin from plugins requiredBy fields
	p.cleanRequiredBy(pluginName)

	delete(p.Plugins, pluginName)
}

// Load plugins warnings into PluginManager struct
func (pm *PluginManager) LoadWarnings() {
	logrus.Debugln("LoadWarnings executed")
	for _, plugin := range pm.Plugins {
		// clear warnings but keep allocated memory
		plugin.Warnings = plugin.Warnings[:0]

		for _, warn := range pm.UpdateCenter.Warnings {

			// skip all plugins except current one
			if warn.Name == plugin.Name {

				for _, warningVersion := range warn.Versions {
					// if warning version lower then current than we dont have warning
					if warningVersion.LastVersion == plugin.Version || utils.IsNewerThan(warningVersion.LastVersion, plugin.Version) {
						// write errors back to plugin
						// TODO: rewrite resolve deps function, may reuse warnings list when newer plugin comes
						plugin.Warnings = append(plugin.Warnings, Warnings{
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
}

func (pm *PluginManager) FixWarnings() {
	for _, plugin := range pm.Plugins {

		if len(plugin.Warnings) == 0 {
			logrus.Debugf("No error found for plugin %s version %s", plugin.Name, plugin.Version)
			// return nil
		}

		var versions []string

		for version := range pm.PluginVersions.Plugins[plugin.Name] {
			versions = append(versions, version)
		}

		// TODO: refact
		var nextVersion string
		var currentVersion string = plugin.Version
		// currentVersion = plugin.Version
		var newPlugin Plugin
		for {
			var err error
			nextVersion, err = utils.GetNextVersion(versions, currentVersion)

			if nextVersion == "" {
				logrus.Infoln("Reach get version limit")
			}

			if err != nil {
				logrus.Warnf("Error during GetNextVersion for %s\n", currentVersion)
				logrus.Warnln(err)
				// return err
			}
			currentVersion = nextVersion

			newPlugin = Plugin{
				Name:    plugin.Name,
				Url:     plugin.Url,
				Version: nextVersion,
			}

			// newPlugin.LoadWarnings()

			fmt.Printf("pl : %s version %s warn: %d", newPlugin.Name, newPlugin.Version, len(newPlugin.Warnings))

			if len(newPlugin.Warnings) == 0 {
				break
			}
		}

		// set new version as current and remove deps as could differs
		plugin.Version = newPlugin.Version
		plugin.Url = newPlugin.Url
		plugin.Warnings = newPlugin.Warnings
		plugin.Dependencies = make(map[string]Plugin)
		// TODO: error if latest version also has error
		// p.predownloadPluginData()

	}
}

func (pm *PluginManager) FixPluginDependenciesMinimal() map[string]Plugin {
	return map[string]Plugin{}
}

// TODO: this function might add new version or new plugin and we dont have data like url + deps etc for it
// Will get plugin dependencies for plugin list and return result back to PluginManager struct
// Will replace version of plugin if dep version > current version
// Used LoadDependenciesFromManifest function as source of versions/plugin so will download hpi files
// WITH CORE UPGRADE
func (pm *PluginManager) FixPluginDependencies() map[string]Plugin {

	pluginsToCheck := make(map[string]Plugin) // plugins has recently updated version or new dependency to check
	pluginsChecked := make(map[string]Plugin) // already checked plugins

	// convert all plugins to pluginToCheck
	for _, pluginPrimary := range pm.Plugins {
		// TODO: SRP
		if pluginPrimary.Url == "" {
			pluginPrimary.Url = pm.PluginVersions.Plugins[pluginPrimary.Name][pluginPrimary.Version].Url
		}

		pluginsToCheck[pluginPrimary.Name] = *pluginPrimary
		logrus.Infof("[FixPluginDependencies] added initial plugin to pluginToCheck %s", pluginPrimary.Name)
	}

	// while we have plugin to check do loop
	for len(pluginsToCheck) > 0 {
		// Iterate over pluginsToCheck
		for _, plugin := range pluginsToCheck {
			logrus.Infof("[FixPluginDependencies] checking plugin: %s:%s\n", plugin.Name, plugin.Version)

			for _, dep := range plugin.Dependencies {
				// plugin.Dependencies[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)

				logrus.Infof("[FixPluginDependencies] checking dep: %s:%s for plugin %s:%s\n", dep.Name, dep.Version, plugin.Name, plugin.Version)

				// if found in pluginsChecked
				if _, found := pluginsChecked[dep.Name]; found {

					if utils.IsNewerThan(dep.Version, pluginsChecked[dep.Name].Version) {
						logrus.Infof("[FixPluginDependencies] Upgrading pluginsChecked bundled dependency %s:%s -> %s:%s\n", pluginsChecked[dep.Name].Name, pluginsChecked[dep.Name].Version, dep.Name, dep.Version)
						// if dependency exist in plugin checked and dep version higher than we have
						pluginsToCheck[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)
						delete(pluginsChecked, dep.Name)
					} else {
						logrus.Infof("[FixPluginDependencies] Skipping pluginsChecked already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsChecked[dep.Name].Version, dep.Version, pluginsChecked[dep.Name].Version)
						// if dependency exist in plugin checked and dep version lower than we have
						// TODO: why we need this?
						// pluginsChecked[dep.Name] = *NewPluginWithVersion(pluginsChecked[dep.Name].Name, pluginsChecked[dep.Name].Version)

					}
					// if found in pluginsToChecked
					// try to find in pluginsToCheck list
				} else if _, foundItem := pluginsToCheck[dep.Name]; foundItem {
					if utils.IsNewerThan(dep.Version, pluginsToCheck[dep.Name].Version) {
						logrus.Infof("[FixPluginDependencies] Upgrading pluginsToCheck bundled dependency %s:%s -> %s:%s\n", pluginsToCheck[dep.Name].Name, pluginsToCheck[dep.Name].Version, dep.Name, dep.Version)
						pluginsToCheck[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)
					} else {
						logrus.Infof("[FixPluginDependencies] Skipping pluginsToCheck already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsToCheck[dep.Name].Version, dep.Version, pluginsToCheck[dep.Name].Version)
						// move plugin from tocheck list to checked list

						// pluginsToCheck[dep.Name].RequiredBy[plugin.Name] = plugin.Version
						// TODO: what that part does? and why we need it?
						// remove already checked plugin
						// NOTE: if found in pluginToCheck and version lower than we have
						// if pluginsToCheck iter var == dependency of itered plugin than move to checked plugin and delete from pluginToCheck
						if plugin.Name == dep.Name {
							pluginsChecked[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)
							logrus.Infoln("[FixPluginDependencies] DELETE PLUGIN FROM pluginsToCheck" + dep.Name)
							delete(pluginsToCheck, dep.Name)
						}
					}
					// didn't find in pluginsChecked and in pluginsToCheck
					// assuming new dependency
				} else {
					// do not load optional deps
					// TODO: make env config to set this
					if dep.Optional {
						logrus.Infof("[FixPluginDependencies] Skipped optional plugin %s", dep.Name)
						continue
					}
					logrus.Infof("[FixPluginDependencies] added new dep to pluginsToCheck %s:%s \n", dep.Name, dep.Version)
					pluginsToCheck[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)
				}
			}

			logrus.Infof("[FixPluginDependencies] added current to pluginsChecked %s:%s\n", plugin.Name, plugin.Version)
			pluginsChecked[plugin.Name] = *NewPluginWithVersion(plugin.Name, plugin.Version)
			delete(pluginsToCheck, plugin.Name)
		}
	}

	// clean what we have in the slice
	pm.UpdatedPlugins = make(map[string]*Plugin)

	for _, pl := range pluginsChecked {
		// create new object and load into UpdatedPlugins
		newpl := NewPluginWithVersion(pl.Name, pl.Version)
		// load plugin data into plugin
		pm.LoadPluginData(newpl)
		pm.UpdatedPlugins[pl.Name] = newpl
	}

	return pluginsChecked
}

func (pm *PluginManager) SetUpdatedPluginWithVersion(pluginName string, pluginVersion string) {
	pm.UpdatedPlugins[pluginName] = NewPluginWithVersion(pluginName, pluginVersion)
}

func (pm *PluginManager) GetUpdatedPlugins() map[string]*Plugin {
	return pm.UpdatedPlugins

}

type pluginChangedType int

const (
	UNKNOWN pluginChangedType = iota
	NEW_PLUGIN
	UPDATED_PLUGIN    // used as dependency for some plugins
	NO_CHANGED_PLUGIN // no one rely on it as dep
)

func (pt pluginChangedType) String() string {
	return []string{"Unknown", "External", "Transitive", "Primary"}[pt]
}

type diffPlugins struct {
	Name           string
	CurrentVersion string
	NewVersion     string
	HTMLURL        string
	ReleaseNotes   []types.ReleaseNote
	// TODO: make enum
	// 1 new 2 update 3 the same
	Type pluginChangedType
}

// Loads dependencies from jenkins update center into Plugin struct
func (pm *PluginManager) LoadDependenciesFromUpdateCenter(p *Plugin) map[string]Plugin {
	for _, dep := range pm.PluginVersions.Plugins[p.Name][p.Version].Dependencies {
		if !dep.Optional {

			p.Dependencies[dep.Name] = Plugin{
				Name:    dep.Name,
				Version: dep.Version,
			}

		}
	}
	return p.Dependencies
}

// Get diff current plugins and updatedPlugins
func (pm *PluginManager) GetFixedDepsDiff() []diffPlugins {
	logrus.Infoln("[GetFixedDepsDiff]")
	var resultDiff []diffPlugins

	for _, changedPlugin := range pm.UpdatedPlugins {
		existingPlugin, exists := pm.Plugins[changedPlugin.Name]

		if exists {
			// if we already have this plugin and version updated
			if utils.IsNewerThan(changedPlugin.Version, existingPlugin.Version) {
				// Get release notes
				// releaseNotes, _ := sources.DownloadProjectReleaseNotes(&pm.PluginSite, changedPlugin.Name)
				releaseNotes, _ := sources.DownloadProjectReleaseNotes(&pm.GitHubClient, changedPlugin.Name)

				var resultRelaseNotes []types.ReleaseNote

				foundNewVersion := false
				foundOldVersion := false
				for _, releaseNote := range releaseNotes {

					if releaseNote.Name == existingPlugin.Version {
						foundOldVersion = true
						continue
					}
					if releaseNote.Name == changedPlugin.Version {
						foundNewVersion = true
					}

					if foundNewVersion || foundOldVersion {
						resultRelaseNotes = append(resultRelaseNotes, releaseNote)
					}

					if foundNewVersion && foundOldVersion {
						break
					}

				}
				if !foundNewVersion {
					logrus.Warnf("haven't found version for diff %s:%s\n", changedPlugin.Name, changedPlugin.Version)
				}

				diff := diffPlugins{
					Name:           changedPlugin.Name,
					CurrentVersion: existingPlugin.Version,
					NewVersion:     changedPlugin.Version,
					HTMLURL:        changedPlugin.GITUrl,
					Type:           UPDATED_PLUGIN,
					ReleaseNotes:   []types.ReleaseNote{},
				}
				if len(resultRelaseNotes) > 0 {
					diff.ReleaseNotes = resultRelaseNotes
				}
				resultDiff = append(resultDiff, diff)
				continue
			}
			// Get release notes end ^
			if existingPlugin.Version == changedPlugin.Version {
				// if version the same
				resultDiff = append(resultDiff, diffPlugins{
					Name:           changedPlugin.Name,
					CurrentVersion: existingPlugin.Version,
					NewVersion:     existingPlugin.Version,
					HTMLURL:        changedPlugin.GITUrl,
					Type:           NO_CHANGED_PLUGIN,
				})
				continue
			}

		} else {
			// new dep
			resultDiff = append(resultDiff, diffPlugins{
				Name:           changedPlugin.Name,
				CurrentVersion: "",
				NewVersion:     changedPlugin.Version,
				HTMLURL:        changedPlugin.GITUrl,
				Type:           NEW_PLUGIN,
			})
			continue
		}
		logrus.Errorln("[GetFixedDepsDiff] SHOULD NOT REACH HERE")
	}

	return resultDiff
}

func (pm *PluginManager) CleanPlugins() {
	// clean what we have in slice
	for k := range pm.Plugins {
		delete(pm.Plugins, k)
	}
}

func (pm *PluginManager) CleanUpdatedPlugins() {
	// clean what we have in slice
	for k := range pm.UpdatedPlugins {
		delete(pm.UpdatedPlugins, k)
	}
}
