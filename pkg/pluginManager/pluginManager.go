package pluginManager

import (
	"fmt"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/updateCenter/pluginVersions"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/updateCenter/updateCenter"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/sirupsen/logrus"
)

/* TODO:
* fix URL
* fix TYPE
* add func to find plugins without dep
* add remove plugin func
* add logic to add plugin
* make dump to file flexible
 */

type PluginManager struct {
	Plugins        map[string]*Plugin // we use map to speedup get set find
	updateCenter   *updateCenter.UpdateCenter
	pluginVersions *pluginVersions.PluginVersions
	coreVersion    string // todo set core version
}

func NewPluginManager() PluginManager {
	uc, _ := updateCenter.Get("")
	pv, _ := pluginVersions.Get()
	return PluginManager{
		Plugins:        make(map[string]*Plugin),
		updateCenter:   uc,
		pluginVersions: pv,
	}
}

func (pm *PluginManager) preloadPluginData(p *Plugin) {
	// Load deps
	for _, dep := range pm.pluginVersions.Plugins[p.Name][p.Version].Dependencies {
		if !dep.Optional {
			p.Dependencies[dep.Name] = Plugin{
				Name:    dep.Name,
				Version: dep.Version,
			}
		}
	}
	// Load url
	p.Url = pm.pluginVersions.Plugins[p.Name][p.Version].Url
	// Load required core version
	p.RequiredCoreVersion = pm.pluginVersions.Plugins[p.Name][p.Version].RequiredCore
}

// TODO: return error?
func (pm *PluginManager) AddPlugin(pl *Plugin) {
	logrus.Infof("Adding new plugin to pluginManager %s:%s", pl.Name, pl.Version)
	if _, found := pm.Plugins[pl.Name]; found {
		logrus.Warnf("Found copy of plugin in pluginsfile. Plugin name: '%s'\n", pl.Name)
	}

	pm.preloadPluginData(pl)

	pm.Plugins[pl.Name] = pl
}

// Load plugins warnings into PluginManager struct
func (pm *PluginManager) LoadWarnings() {
	for _, plugin := range pm.Plugins {

		logrus.Debugln("LoadWarnings executed")
		// download plugin hpi file
		// p.Download()

		// requiredCoreVersion := p.GetRequiredCoreVersion()
		// uc, _ := updateCenter.Get(p.RequiredCoreVersion)

		// clear warnings but keep allocated memory
		plugin.Warnings = plugin.Warnings[:0]

		for _, warn := range pm.updateCenter.Warnings {

			// skip all plugins except current one
			if warn.Name == plugin.Name {

				// fmt.Println(warn.Message)
				// fmt.Println(warn.Versions)
				// fmt.Println(warn.Id)

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

		// TODO: do not call each time reload warnings everywhere
		// p.LoadWarnings() // we need to find last version with warning
		if len(plugin.Warnings) == 0 {
			logrus.Debugf("No error found for plugin %s version %s", plugin.Name, plugin.Version)
			// return nil
		}
		// pv, _ := pluginVersions.Get() // to check if there any version where warning version +1
		logrus.Infoln(pm.pluginVersions.Plugins[plugin.Name][plugin.Version].RequiredCore)

		var versions []string

		for version := range pm.pluginVersions.Plugins[plugin.Name] {
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

// func (pm *PluginManager) FixPluginWarnings(p *Plugin) {

// 	// clear warnings but keep allocated memory
// 	p.Warnings = p.Warnings[:0]

// 	for _, warn := range pm.updateCenter.Warnings {

// 		// skip all plugins except current one
// 		if warn.Name == p.Name {

// 			for _, warningVersion := range warn.Versions {
// 				// if warning version lower then current than we dont have warning
// 				if warningVersion.LastVersion == p.Version || utils.IsNewerThan(warningVersion.LastVersion, p.Version) {
// 					// write errors back to plugin
// 					// TODO: rewrite resolve deps function, may reuse warnings list when newer plugin comes
// 					p.Warnings = append(p.Warnings, Warnings{
// 						Id:      warn.Id,
// 						Message: warn.Message,
// 						Name:    warn.Name,
// 						Url:     warn.Url,
// 						Versions: []struct {
// 							LastVersion string
// 							Pattern     string
// 						}(warn.Versions),
// 					})
// 				}
// 			}
// 		}
// 	}
// }

/*
	PluginA has dep PluginB
	PluginC has no deps
	PluginD has dep PluginA

	PluginC - doesnt have deps so we can easily delete this plugin
	PluginD - Primary type because no one rely on it

*/
// find plugins plugins which no one rely on
func (p *PluginManager) LoadPluginTypes() {}

// will try to delete plugin from list if no one rely as dep on this plugin
func (p *PluginManager) DeletePlugin(pluginName string) error {
	// check if we can delete plugin and no one rely on it
	// if some one use it print
	return nil
}

// Predownload PluginManager plugins
// func (p *PluginManager) PredownloadPlugins() {

// 	var wg sync.WaitGroup

// 	for _, plugin := range *p {
// 		wg.Add(1)

// 		go func(plugin *Plugin) {
// 			defer wg.Done()
// 			plugin.Download()
// 		}(plugin)

// 	}
// 	wg.Wait()
// }

// TODO: add tests
// get plugins which no one rely on
func (pm *PluginManager) GetStandalonePlugins() []*Plugin {
	var plugins []*Plugin
	for _, plugin := range pm.Plugins {
		if len(plugin.RequiredBy) == 0 {

			plugins = append(plugins, plugin)
		}
	}

	return plugins
}

// TODO: this function might add new version or new plugin and we dont have data like url + deps etc for it
// Will get plugin dependencies for plugin list and return result back to PluginManager struct
// Will replace version of plugin if dep version > current version
// Used LoadDependenciesFromManifest function as source of versions/plugin so will download hpi files
func (pm *PluginManager) FixPluginDependencies() {

	pluginsToCheck := make(map[string]Plugin) // plugins has recently updated version or new dependency to check
	pluginsChecked := make(map[string]Plugin) // already checked plugins

	// convert all plugins to pluginToCheck
	for _, pluginPrimary := range pm.Plugins {
		if pluginPrimary.Url == "" {
			pluginPrimary.Url = pm.pluginVersions.Plugins[pluginPrimary.Name][pluginPrimary.Version].Url
		}

		pluginsToCheck[pluginPrimary.Name] = *pluginPrimary
		logrus.Infof("added initial plugin to pluginToCheck %s", pluginPrimary.Name)
	}

	// while we have plugin to check do loop
	for len(pluginsToCheck) > 0 {
		// Iterate over pluginsToCheck
		for _, plugin := range pluginsToCheck {
			logrus.Infof("checking plugin: %s:%s\n", plugin.Name, plugin.Version)

			for _, dep := range plugin.LoadDependenciesFromUpdateCenter() {
				plugin.Dependencies[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)

				logrus.Infof("checking dep: %s:%s for plugin %s:%s\n", dep.Name, dep.Version, plugin.Name, plugin.Version)

				// if found in pluginsChecked
				if _, found := pluginsChecked[dep.Name]; found {

					if utils.IsNewerThan(dep.Version, pluginsChecked[dep.Name].Version) {
						logrus.Infof("Upgrading pluginsChecked bundled dependency %s:%s -> %s:%s\n", pluginsChecked[dep.Name].Name, pluginsChecked[dep.Name].Version, dep.Name, dep.Version)
						// if dependency exist in plugin checked and dep version higher than we have

						pluginsChecked[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						pluginsToCheck[dep.Name] = Plugin{
							Name:         dep.Name,
							Version:      dep.Version,
							Url:          dep.Url,
							Type:         pluginsChecked[dep.Name].Type,
							Dependencies: pluginsChecked[dep.Name].Dependencies,
							RequiredBy:   pluginsChecked[dep.Name].RequiredBy,
						}
						delete(pluginsChecked, dep.Name)
					} else {
						logrus.Infof("Skipping pluginsChecked already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsChecked[dep.Name].Version, dep.Version, pluginsChecked[dep.Name].Version)
						// if dependency exist in plugin checked and dep version lower than we have

						pluginsChecked[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						// TODO: why we need this?
						pluginsChecked[dep.Name] = Plugin{
							Name:         pluginsChecked[dep.Name].Name,
							Version:      pluginsChecked[dep.Name].Version,
							Type:         pluginsChecked[dep.Name].Type,
							Url:          pluginsChecked[dep.Name].Url,
							Dependencies: pluginsChecked[dep.Name].Dependencies,
							RequiredBy:   pluginsChecked[dep.Name].RequiredBy,
						}

					}
					// if found in pluginsToChecked
					// try to find in pluginsToCheck list
				} else if _, foundItem := pluginsToCheck[dep.Name]; foundItem {
					if utils.IsNewerThan(dep.Version, pluginsToCheck[dep.Name].Version) {
						logrus.Infof("Upgrading pluginsToCheck bundled dependency %s:%s -> %s:%s\n", pluginsToCheck[dep.Name].Name, pluginsToCheck[dep.Name].Version, dep.Name, dep.Version)

						pluginsToCheck[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						pluginsToCheck[dep.Name] = Plugin{
							Name:         dep.Name,
							Version:      dep.Version,
							Type:         plugin.Type,
							Url:          pm.pluginVersions.Plugins[dep.Name][dep.Version].Url,
							Dependencies: pluginsToCheck[dep.Name].Dependencies,
							RequiredBy:   pluginsToCheck[dep.Name].RequiredBy,
						}
					} else {
						logrus.Infof("Skipping pluginsToCheck already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsToCheck[dep.Name].Version, dep.Version, pluginsToCheck[dep.Name].Version)
						// move plugin from tocheck list to checked list

						pluginsToCheck[dep.Name].RequiredBy[plugin.Name] = plugin.Version
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
								RequiredBy:   pluginsToCheck[dep.Name].RequiredBy,
							}
							logrus.Infoln("DELETE PLUGIN FROM pluginsToCheck" + dep.Name)
							delete(pluginsToCheck, dep.Name)
						}
					}
					// didn't find in pluginsChecked and in pluginsToCheck
					// assuming new dependency
				} else {
					// do not load optional deps
					// TODO: make env config to set this
					if dep.Optional {
						logrus.Infof("Skipped optional plugin %s", dep.Name)
						continue
					}
					logrus.Infof("added new dep to pluginsToCheck %s:%s \n", dep.Name, dep.Version)
					pluginsToCheck[dep.Name] = Plugin{
						Name:         dep.Name,
						Version:      dep.Version,
						Type:         UNKNOWN,
						Url:          "55555",
						Dependencies: make(map[string]Plugin),
						RequiredBy: map[string]string{
							plugin.Name: plugin.Version,
						},
					}
				}
			}

			logrus.Infof("added current to pluginsChecked %s:%s\n", plugin.Name, plugin.Version)
			pluginsChecked[plugin.Name] = Plugin{
				Name:         plugin.Name,
				Version:      plugin.Version,
				Type:         plugin.Type,
				Dependencies: plugin.Dependencies,
				Url:          plugin.Url,
				RequiredBy:   plugin.RequiredBy,
			}

			delete(pluginsToCheck, plugin.Name)
		}
	}

	// clean what we have in slice
	for k := range pm.Plugins {
		delete(pm.Plugins, k)
	}

	// write back to PluginManager struct
	for _, pl := range pluginsChecked {
		pm.Plugins[pl.Name] = &Plugin{
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
}

/*
TODO: require jenkins core version?
go over plugins with errors and try to update plugin to new version which might not have warn
requires FixPluginDependencies after list patch and run FixPluginWarnings again?
TODO: find better logic then recursion ^
like check all plugin version to latest and check warnings
need to be aware might require jenkins core version patch
*/

// TODO: remove unused
func (p *PluginManager) GeneratePluginDepsFile() {
	fmt.Println("NOT IMPLEMENTED")
}

// TODO: remove unused
func (pm *PluginManager) GetLatestVersions(uc *updateCenter.UpdateCenter) {

	for _, plugins := range pm.Plugins {
		plugins.Version = uc.Plugins[plugins.Name].Version
	}

	// p.UpdateFile(p.File.fileName)
}
