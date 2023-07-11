package pluginManager

import (
	"fmt"
	"sort"

	"github.com/emelianrus/jenkins-release-notes-parser/outputGenerators"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/parsers"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/sources"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/sources/github"
	jenkins "github.com/emelianrus/jenkins-release-notes-parser/pkg/sources/jenkinsPluginSite"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/updateCenter/pluginVersions"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/updateCenter/updateCenter"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/sirupsen/logrus"
)

type PluginManager struct {
	Plugins map[string]*Plugin // Main list where we have plugins. We use map to speedup get set find

	UpdatedPlugins map[string]*Plugin // Temp storage. will set this after update/fix deps

	coreVersion string // jenkins core version

	UpdateCenter   *updateCenter.UpdateCenter     // external. from jenkins api. plugin latest version + deprecations + deps for plugin
	PluginVersions *pluginVersions.PluginVersions // external. from jenkins api. plugins data with versions

	// release notes sources
	PluginSite   jenkins.PluginSite // jenkins site which has release notes (last 10)
	GitHubClient github.GitHub      // github client to get release notes

	// file read/write logic
	FileParser parsers.InputParser            // parses txt file into plugin-manager plugins
	FileOutput outputGenerators.FileGenerator // plugin-manager plugins into file content
}

func NewPluginManager() PluginManager {
	uc, _ := updateCenter.Get("")
	pv, _ := pluginVersions.Get()

	return PluginManager{
		coreVersion: "2.235.2", // TODO: should not be hardcoded

		Plugins:        make(map[string]*Plugin),
		UpdatedPlugins: make(map[string]*Plugin),

		UpdateCenter:   uc,
		PluginVersions: pv,

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

func (pm *PluginManager) SetFileOutput(og outputGenerators.FileGenerator) {
	pm.FileOutput = outputGenerators.SetOutputGenerator(og)
}

func (pm *PluginManager) GenerateFileOutput() []byte {
	result := make(map[string]string)

	for _, pl := range pm.UpdatedPlugins {
		result[pl.Name] = pl.Version
	}
	return pm.FileOutput.Generate(result)

}

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

func (pm *PluginManager) preloadPluginData(p *Plugin) {
	logrus.Infof("preloadPluginData %s:%s", p.Name, p.Version)
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

func (pm *PluginManager) GetCoreVersion() string {
	return pm.coreVersion
}
func (pm *PluginManager) SetCoreVersion(newCoreVersion string) {
	pm.coreVersion = newCoreVersion
	uc, _ := updateCenter.Get(newCoreVersion)

	pm.UpdateCenter = uc
}

func (pm *PluginManager) AddPluginWithVersion(pluginName string, version string) {
	logrus.Infof("[AddPluginWithVersion] Adding new plugin to pluginManager %s:%s", pluginName, version)
	pl := NewPluginWithVersion(pluginName, version)

	logrus.Infof("Adding new plugin to pluginManager %s:%s", pl.Name, pl.Version)
	if _, found := pm.Plugins[pl.Name]; found {
		logrus.Warnf("Found copy of plugin in pluginsfile. Plugin name: '%s'\n", pl.Name)
	}

	pm.preloadPluginData(pl)
	// TODO: sync with DB

	pm.Plugins[pl.Name] = pl
	pm.generateRequiredBy()
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

// TODO: this function might add new version or new plugin and we dont have data like url + deps etc for it
// Will get plugin dependencies for plugin list and return result back to PluginManager struct
// Will replace version of plugin if dep version > current version
// Used LoadDependenciesFromManifest function as source of versions/plugin so will download hpi files
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

			for _, dep := range plugin.LoadDependenciesFromUpdateCenter() {
				// plugin.Dependencies[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)

				logrus.Infof("[FixPluginDependencies] checking dep: %s:%s for plugin %s:%s\n", dep.Name, dep.Version, plugin.Name, plugin.Version)

				// if found in pluginsChecked
				if _, found := pluginsChecked[dep.Name]; found {

					if utils.IsNewerThan(dep.Version, pluginsChecked[dep.Name].Version) {
						logrus.Infof("[FixPluginDependencies] Upgrading pluginsChecked bundled dependency %s:%s -> %s:%s\n", pluginsChecked[dep.Name].Name, pluginsChecked[dep.Name].Version, dep.Name, dep.Version)
						// if dependency exist in plugin checked and dep version higher than we have

						// pluginsChecked[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						pluginsToCheck[dep.Name] = Plugin{
							Name:         dep.Name,
							Version:      dep.Version,
							Url:          pm.PluginVersions.Plugins[dep.Name][dep.Version].Url,
							Type:         pluginsChecked[dep.Name].Type,
							Dependencies: pluginsChecked[dep.Name].Dependencies,
							// RequiredBy:   pluginsChecked[dep.Name].RequiredBy,
						}
						delete(pluginsChecked, dep.Name)
					} else {
						logrus.Infof("[FixPluginDependencies] Skipping pluginsChecked already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsChecked[dep.Name].Version, dep.Version, pluginsChecked[dep.Name].Version)
						// if dependency exist in plugin checked and dep version lower than we have

						// pluginsChecked[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						// TODO: why we need this?
						pluginsChecked[dep.Name] = Plugin{
							Name:         pluginsChecked[dep.Name].Name,
							Version:      pluginsChecked[dep.Name].Version,
							Type:         pluginsChecked[dep.Name].Type,
							Url:          pluginsChecked[dep.Name].Url,
							Dependencies: pluginsChecked[dep.Name].Dependencies,
							// RequiredBy:   pluginsChecked[dep.Name].RequiredBy,
						}

					}
					// if found in pluginsToChecked
					// try to find in pluginsToCheck list
				} else if _, foundItem := pluginsToCheck[dep.Name]; foundItem {
					if utils.IsNewerThan(dep.Version, pluginsToCheck[dep.Name].Version) {
						logrus.Infof("[FixPluginDependencies] Upgrading pluginsToCheck bundled dependency %s:%s -> %s:%s\n", pluginsToCheck[dep.Name].Name, pluginsToCheck[dep.Name].Version, dep.Name, dep.Version)

						// pluginsToCheck[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						pluginsToCheck[dep.Name] = Plugin{
							Name:         dep.Name,
							Version:      dep.Version,
							Type:         plugin.Type,
							Url:          pm.PluginVersions.Plugins[dep.Name][dep.Version].Url,
							Dependencies: pluginsToCheck[dep.Name].Dependencies,
							// RequiredBy:   pluginsToCheck[dep.Name].RequiredBy,
						}
					} else {
						logrus.Infof("[FixPluginDependencies] Skipping pluginsToCheck already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsToCheck[dep.Name].Version, dep.Version, pluginsToCheck[dep.Name].Version)
						// move plugin from tocheck list to checked list

						// pluginsToCheck[dep.Name].RequiredBy[plugin.Name] = plugin.Version
						// TODO: what that part does? and why we need it?
						// remove already checked plugin
						// NOTE: if found in pluginToCheck and version lower than we have
						// if pluginsToCheck iter var == dependency of itered plugin than move to checked plugin and delete from pluginToCheck
						if plugin.Name == dep.Name {
							pluginsChecked[dep.Name] = Plugin{
								Name:         dep.Name,
								Version:      dep.Version,
								Type:         plugin.Type,
								Url:          "4444",
								Dependencies: pluginsToCheck[dep.Name].Dependencies,
								// RequiredBy:   pluginsToCheck[dep.Name].RequiredBy,
							}
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
					pluginsToCheck[dep.Name] = Plugin{
						Name:         dep.Name,
						Version:      dep.Version,
						Type:         UNKNOWN,
						Url:          pm.PluginVersions.Plugins[dep.Name][dep.Version].Url,
						Dependencies: make(map[string]Plugin),
						RequiredBy: map[string]string{
							plugin.Name: plugin.Version,
						},
					}
				}
			}

			logrus.Infof("[FixPluginDependencies] added current to pluginsChecked %s:%s\n", plugin.Name, plugin.Version)
			pluginsChecked[plugin.Name] = Plugin{
				Name:         plugin.Name,
				Version:      plugin.Version,
				Type:         plugin.Type,
				Dependencies: plugin.Dependencies,
				Url:          plugin.Url,
				// RequiredBy:   plugin.RequiredBy,
			}

			delete(pluginsToCheck, plugin.Name)
		}
	}

	// clean what we have in slice
	for k := range pm.UpdatedPlugins {
		delete(pm.UpdatedPlugins, k)
	}

	for _, pl := range pluginsChecked {
		pm.UpdatedPlugins[pl.Name] = &Plugin{
			Name:         pl.Name,
			Version:      pl.Version,
			Url:          pl.Url,
			Dependencies: pl.Dependencies,
			RequiredBy:   pl.RequiredBy,
			Optional:     pl.Optional,
			Warnings:     pl.Warnings,
			Type:         pl.Type,
		}
	}

	return pluginsChecked
}

func (pm *PluginManager) SetUpdatedPluginWithVersion(pluginName string, pluginVersion string) {
	pm.UpdatedPlugins[pluginName] = NewPluginWithVersion(pluginName, pluginVersion)
}

func (pm *PluginManager) GetUpdatedPlugins() map[string]*Plugin {
	return pm.UpdatedPlugins
}

type diffPlugins struct {
	Name           string
	CurrentVersion string
	NewVersion     string
	ReleaseNotes   []types.ReleaseNote
	// TODO: make enum
	// 1 new 2 update 3 the same
	Type int
}

// Get diff current plugins and updatedPlugins
func (pm *PluginManager) GetFixedDepsDiff() []diffPlugins {
	logrus.Infoln("[GetFixedDepsDiff]")
	var resultDiff []diffPlugins

	for _, changedPlugin := range pm.UpdatedPlugins {
		if _, exists := pm.Plugins[changedPlugin.Name]; exists {
			// if we already have this plugin and version updated
			if utils.IsNewerThan(changedPlugin.Version, pm.Plugins[changedPlugin.Name].Version) {
				// Get release notes

				// releaseNotes, _ := sources.DownloadProjectReleaseNotes(&pm.PluginSite, changedPlugin.Name)
				gh := github.NewGitHubClient()
				releaseNotes, _ := sources.DownloadProjectReleaseNotes(&gh, changedPlugin.Name)

				var resultRelaseNotes []types.ReleaseNote

				sort.Slice(releaseNotes, func(i, j int) bool {
					return utils.IsNewerThan(releaseNotes[i].Name, releaseNotes[j].Name)
				})

				foundNewVersion := false
				foundOldVersion := false
				for _, releaseNote := range releaseNotes {
					if releaseNote.Name == pm.Plugins[changedPlugin.Name].Version {
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
					logrus.Warnf("haven't found version %s:%s\n", changedPlugin.Name, changedPlugin.Version)
					// TODO: try to download from github
					// add feed
				}

				diff := diffPlugins{
					Name:           changedPlugin.Name,
					CurrentVersion: pm.Plugins[changedPlugin.Name].Version,
					NewVersion:     changedPlugin.Version,
					Type:           2, // update
					ReleaseNotes:   []types.ReleaseNote{},
				}
				if len(resultRelaseNotes) > 0 {
					diff.ReleaseNotes = resultRelaseNotes
				}
				resultDiff = append(resultDiff, diff)
				continue
			}
			// Get release notes end ^
			if pm.Plugins[changedPlugin.Name].Version == changedPlugin.Version {
				// if version the same
				resultDiff = append(resultDiff, diffPlugins{
					Name:           changedPlugin.Name,
					CurrentVersion: pm.Plugins[changedPlugin.Name].Version,
					NewVersion:     pm.Plugins[changedPlugin.Name].Version,
					Type:           3, // version the same
				})
				continue
			}

		} else {
			// new dep
			resultDiff = append(resultDiff, diffPlugins{
				Name:           changedPlugin.Name,
				CurrentVersion: "",
				NewVersion:     changedPlugin.Version,
				Type:           1, // new plugin
			})
			continue
		}
		logrus.Errorln("[GetFixedDepsDiff] SHOULD NOT REACH HERE")
	}

	return resultDiff
}
