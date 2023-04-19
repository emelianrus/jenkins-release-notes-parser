package pluginManager

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"

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

// TODO: make this[]Plugin???
// we use map to speedup get set find
type PluginManager map[string]*Plugin

func NewPluginManager() PluginManager {
	return PluginManager{}
}

// TODO: return error?
func (pm *PluginManager) AddPlugin(pl *Plugin) {
	logrus.Infof("Adding new plugin to pluginManager %s:%s", pl.Name, pl.Version)
	if _, found := (*pm)[pl.Name]; found {
		logrus.Warnf("Found copy of plugin in pluginsfile. Plugin name: '%s'\n", pl.Name)
	}
	(*pm)[pl.Name] = pl
}

// Load plugins warnings into PluginManager struct
func (p *PluginManager) LoadWarnings() {
	for _, plugin := range *p {
		plugin.LoadWarnings()
	}
}

// TODO: require jenkins core version?
// go over plugins with errors and try to update plugin to new version which might not have warn
// requires FixPluginDependencies after list patch and run FixPluginWarnings again?
// TODO: find better logic than recursion ^
// like check all plugin version to latest and check warnings
// need to be aware might require jenkins core version patch
func (p *PluginManager) FixWarnings() {
	for _, plugin := range *p {
		plugin.FixWarnings()
		plugin.LoadRequiredCoreVersion()
	}

	p.FixPluginDependencies()
}

func (p *PluginManager) DumpToFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	// Sort items by alph
	keys := make([]string, 0, len(*p))
	for k := range *p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		w.WriteString((*p)[k].Name + ":" + (*p)[k].Version + "\n")
	}

	w.Flush()
}

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
func (p *PluginManager) PredownloadPlugins() {

	var wg sync.WaitGroup

	for _, plugin := range *p {
		wg.Add(1)

		go func(plugin *Plugin) {
			defer wg.Done()
			plugin.Download()
		}(plugin)

	}
	wg.Wait()
}

// TODO: add tests
// get plugins which no one rely on
func (p *PluginManager) GetStandalonePlugins() []*Plugin {
	var plugins []*Plugin
	for _, plugin := range *p {
		if len(plugin.RequiredBy) == 0 {
			plugins = append(plugins, plugin)
		}
	}

	return plugins
}

// Will get plugin dependencies for plugin list and return result back to PluginManager struct
// Will replace version of plugin if dep version > current version
// Used LoadDependenciesFromManifest function as source of versions/plugin so will download hpi files
func (p *PluginManager) FixPluginDependencies() {

	pluginsToCheck := make(map[string]Plugin) // plugins has recently updated version or new dependency to check
	pluginsChecked := make(map[string]Plugin) // checked plugins
	// convert all plugins to pluginToCheck
	for _, pluginPrimary := range *p {
		pluginsToCheck[pluginPrimary.Name] = *pluginPrimary
		logrus.Infof("added initial plugin to pluginToCheck %s", pluginPrimary.Name)
	}

	for ok := true; ok; ok = len(pluginsToCheck) > 0 {
		// Iterate over pluginsToCheck
		for _, plugin := range pluginsToCheck {
			logrus.Infof("checking plugin: %s:%s  \n", plugin.Name, plugin.Version)
			// TODO: change to updatecenter + add cache
			for _, dep := range plugin.LoadDependenciesFromUpdateCenter() {
				plugin.Dependencies[dep.Name] = *NewPluginWithVersion(dep.Name, dep.Version)

				logrus.Infof("checking dep: %s:%s for plugin %s:%s\n", dep.Name, dep.Version, plugin.Name, plugin.Version)

				// if found in pluginsChecked
				if _, found := pluginsChecked[dep.Name]; found {

					if utils.IsNewerThan(dep.Version, pluginsChecked[dep.Name].Version) {
						logrus.Infof("Upgrading pluginsChecked bundled dependency %s:%s -> %s:%s\n", pluginsChecked[dep.Name].Name, pluginsChecked[dep.Name].Version, dep.Name, dep.Version)

						pluginsChecked[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						pluginsToCheck[dep.Name] = Plugin{
							Name:         dep.Name,
							Version:      dep.Version,
							Url:          "4444",
							Type:         pluginsChecked[dep.Name].Type,
							Dependencies: pluginsChecked[dep.Name].Dependencies,
							RequiredBy:   pluginsChecked[dep.Name].RequiredBy,
						}
						delete(pluginsChecked, dep.Name)
					} else {
						logrus.Infof("Skipping pluginsChecked already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsChecked[dep.Name].Version, dep.Version, pluginsChecked[dep.Name].Version)

						pluginsChecked[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						pluginsChecked[dep.Name] = Plugin{
							Name:         pluginsChecked[dep.Name].Name,
							Version:      pluginsChecked[dep.Name].Version,
							Type:         pluginsChecked[dep.Name].Type,
							Url:          "3333",
							Dependencies: pluginsChecked[dep.Name].Dependencies,
							RequiredBy:   pluginsChecked[dep.Name].RequiredBy,
						}

					}
					// try to find in pluginsToCheck list
				} else if _, foundItem := pluginsToCheck[dep.Name]; foundItem {
					if utils.IsNewerThan(dep.Version, pluginsToCheck[dep.Name].Version) {
						logrus.Infof("Upgrading pluginsToCheck bundled dependency %s:%s -> %s:%s\n", pluginsToCheck[dep.Name].Name, pluginsToCheck[dep.Name].Version, dep.Name, dep.Version)

						pluginsToCheck[dep.Name].RequiredBy[plugin.Name] = plugin.Version

						pluginsToCheck[dep.Name] = Plugin{
							Name:         dep.Name,
							Version:      dep.Version,
							Type:         plugin.Type,
							Url:          "2222",
							Dependencies: pluginsToCheck[dep.Name].Dependencies,
							RequiredBy:   pluginsToCheck[dep.Name].RequiredBy,
						}
					} else {
						// move plugin from tocheck list to checked list
						logrus.Infof("Skipping pluginsToCheck already installed dependency %s:%s (%s <= %s) \n", dep.Name, pluginsToCheck[dep.Name].Version, dep.Version, pluginsToCheck[dep.Name].Version)

						pluginsToCheck[dep.Name].RequiredBy[plugin.Name] = plugin.Version
						// TODO: what that part does?
						// remove already checked plugin
						if plugin.Name == dep.Name {
							pluginsChecked[dep.Name] = Plugin{
								Name:         dep.Name,
								Version:      dep.Version,
								Type:         plugin.Type,
								Url:          "1111",
								Dependencies: pluginsToCheck[dep.Name].Dependencies,
								RequiredBy:   pluginsToCheck[dep.Name].RequiredBy,
							}
							logrus.Infoln("DELETE PLUGIN FROM pluginsToCheck" + dep.Name)
							delete(pluginsToCheck, dep.Name)
						}
					}
					// dont have in pluginsChecked and in pluginsToCheck
					// assumed new dependency
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
						Url:          dep.Url,
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
	for k := range *p {
		delete(*p, k)
	}

	// write back to PluginManager struct
	for _, pl := range pluginsChecked {
		(*p)[pl.Name] = &Plugin{
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
func (p *PluginManager) GetLatestVersions(uc *updateCenter.UpdateCenter) {

	for _, plugins := range *p {
		plugins.Version = uc.Plugins[plugins.Name].Version
	}

	// p.UpdateFile(p.File.fileName)
}

// TODO: part of fancy output
// newPlugins := make(map[string]Plugin)
// updatedPlugins := make(map[string]Plugin)

// for _, pluginT := range pluginsChecked {
// 	logrus.Infof("WORKING ON PLUGIN %s\n", pluginT.Name)
// 	//logrus.Debugf("%s:%s required by: %s\n", pluginT.Name, pluginT.Version, pluginT.RequiredBy)
// 	if _, found := p.Plugins[pluginT.Name]; found {
// 		plugin := Plugin{
// 			Name:         pluginT.Name,
// 			Version:      pluginT.Version,
// 			Type:         pluginT.Type,
// 			Dependencies: pluginT.Dependencies,
// 			RequiredBy:   pluginT.RequiredBy}

// 		if utils.IsNewerThan(pluginT.Version, p.Plugins[pluginT.Name].Version) {
// 			updatedPlugins[pluginT.Name] = plugin
// 		}

// 	} else {
// 		plugin := Plugin{
// 			Name:         pluginT.Name,
// 			Version:      pluginT.Version,
// 			Dependencies: pluginT.Dependencies,
// 			RequiredBy:   pluginT.RequiredBy,
// 		}
// 		newPlugins[pluginT.Name] = plugin
// 	}
// }
// fmt.Printf("\n\n\n")

// if len(updatedPlugins) > 0 {
// 	fmt.Println("Updated plugin")
// 	for _, updated := range updatedPlugins {
// 		fmt.Printf("%s:%s\n", updated.Name, updated.Version)
// 		for _, p := range updated.RequiredBy {
// 			fmt.Printf("    ↳ required by: %s:%s\n", p.Name, p.Version)
// 		}
// 	}
// }
// if len(newPlugins) > 0 {
// 	fmt.Println("new ")
// 	for _, new := range newPlugins {
// 		fmt.Printf("--- %s:%s\n", new.Name, new.Version)
// 		for _, p := range new.RequiredBy {
// 			fmt.Printf("    ↳ required by: %s:%s\n", p.Name, p.Version)
// 		}
// 	}
// }

// if len(newPlugins) < 1 && len(updatedPlugins) < 1 {
// 	fmt.Println("nothing to do")
// }

// fmt.Println(pluginsChecked)
