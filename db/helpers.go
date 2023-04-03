package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/github"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
)

// DB part end ^

func (r *Redis) GetJenkinsServers() []types.JenkinsServer {
	var servers []types.JenkinsServer

	keys, _ := r.Keys("servers:*")
	for _, path := range keys {
		re := strings.Split(path, ":")
		if len(re) == 2 {

			serverJson, _ := r.Get(path).Bytes()
			var jenkinsServer types.JenkinsServer
			err := json.Unmarshal(serverJson, &jenkinsServer)
			if err != nil {
				fmt.Println("can not unmarshal jenkins server")
			}

			plugins, _ := r.GetJenkinsPlugins(jenkinsServer.Name)
			servers = append(servers, types.JenkinsServer{
				Name:    jenkinsServer.Name,
				Core:    jenkinsServer.Core,
				Plugins: plugins,
			})
		}
	}

	return servers
}

func (r *Redis) AddJenkinsServer(serverName string, coreVersion string) {
	// TODO: check if exist/replace
	// js := JenkinsServer{
	// 	Name: "jenkins-one",
	// 	Core: "2.3233.2",
	// 	Plugins: []ServerPlugin{
	// 		{
	// 			Name:    "plugin-installation-manager-tool",
	// 			Version: "2.10.0",
	// 		},
	// 		{
	// 			Name:    "okhttp-api-plugin",
	// 			Version: "4.9.3-108.v0feda04578cf",
	// 		},
	// 	},
	// }

	jsonData, err := json.Marshal(types.JenkinsServer{
		Name: serverName,
		Core: coreVersion,
	})

	if err != nil {
		log.Println(err)
		return
	}
	// write jenkins server json
	err = r.Set(fmt.Sprintf("servers:%s", serverName), jsonData)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Added %s:%s\n", serverName, coreVersion)
}

func (r *Redis) DeleteJenkinsServer(serverName string) {
	path := fmt.Sprintf("servers:%s", serverName)
	fmt.Printf("Removing jenkins server %s\n", path)
	r.Del(path)
}

func (r *Redis) GetAllPluginsFromServers() []string {
	var result []string
	pluginKeys, _ := r.Keys("servers:*:plugins:*")

	for _, path := range pluginKeys {
		splitted := strings.Split(path, ":")
		result = append(result, splitted[len(splitted)-1])
	}
	result = utils.RemoveDuplicates(result)
	return result
}

func (r *Redis) ChangeJenkinServerPluginVersion(serverName string, pluginName string, newVersion string) error {
	fmt.Printf("Change plugin version to server name: %s plugin name %s new verion %s\n", serverName, pluginName, newVersion)
	jsonData, _ := json.Marshal(newVersion)

	err := r.Set(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName), jsonData)
	if err != nil {
		fmt.Println("can not append to redis new plugin")
	}
	fmt.Printf("version for plugin %s was changed and saved to redis\n", pluginName)

	return nil
}

func (r *Redis) GetJenkinsPlugins(jenkinsServer string) ([]types.JenkinsPlugin, error) {
	var jenkinsPlugins []types.JenkinsPlugin
	pluginKeys, _ := r.Keys(fmt.Sprintf("servers:%s:plugins:*", jenkinsServer))

	for _, pluginKey := range pluginKeys {

		// get plugin name from path (is there is a better way? )
		path := strings.Split(pluginKey, ":")
		pluginName := path[len(path)-1]
		pluginVersionByte, _ := r.Get(pluginKey).Bytes()

		var version string
		err := json.Unmarshal(pluginVersionByte, &version)
		if err != nil {
			fmt.Println("can not unmarshal version")
		}
		projectError := r.GetProjectError(pluginName)

		jenkinsPlugins = append(jenkinsPlugins, types.JenkinsPlugin{
			Name:         pluginName,
			Version:      version,
			Error:        projectError,
			IsDownloaded: r.IsProjectDownloaded(pluginName),
		})
	}

	return jenkinsPlugins, nil
}

func (r *Redis) AddJenkinsServerPlugin(serverName string, plugin types.JenkinsPlugin) error {
	_, err := r.Get(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name)).Bytes()
	if err != nil {
		fmt.Println(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name) + " not found in redis, adding...")

		jsonData, _ := json.Marshal(plugin.Version)

		err := r.Set(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name), jsonData)
		if err != nil {
			fmt.Println("can not append to redis new plugin")
		}
		fmt.Printf("appended plugin %s to redis\n", plugin.Name)
	}
	return nil
}

func (r *Redis) RemoveJenkinsServerPlugin(serverName string, pluginName string) {
	fmt.Printf("removing key from redis %s\n", pluginName)
	r.Del(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName))
}

func (r *Redis) SetLastUpdatedTime(pluginName string, value string) error {

	jsonData, _ := json.Marshal(value)
	err := r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "lastUpdated"),
		jsonData)
	if err != nil {
		log.Println("SetLastUpdatedTime error")
		log.Println(err)
		return nil
	}
	return nil
}

func (r *Redis) GetLastUpdatedTime(pluginName string) string {

	serverJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "lastUpdated")).Bytes()
	return string(serverJson)
}

func (r *Redis) SetProjectError(pluginName string, value string) error {
	jsonData, _ := json.Marshal(value)
	err := r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "error"),
		jsonData)
	if err != nil {
		log.Println("SetProjectError error:")
		log.Println(err)
		return nil
	}
	return nil
}

func (r *Redis) GetProjectError(pluginName string) string {
	serverJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "error")).Bytes()
	return string(serverJson)
}

func (r *Redis) IsProjectDownloaded(pluginName string) bool {
	_, err := r.Get(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "versions")).Bytes()
	if err == nil {
		return true
	} else {
		return false
	}
}

// func (r *Redis) SetPluginWithVersion(pluginName string, pluginVersion string, releaseNote types.GitHubReleaseNote) error {
// 	if pluginName == "" || pluginVersion == "" {
// 		return fmt.Errorf("can not set to DB, name: %s or Version: %s is empty", pluginName, pluginVersion)
// 	}
// 	key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, pluginVersion)
// 	// 0 time.Hour
// 	jsonData, err := json.Marshal(releaseNote)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	err = r.Set(key, jsonData)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	return nil
// }

func (r *Redis) GetPluginWithVersion(pluginName string, pluginVersion string) (types.GitHubReleaseNote, error) {
	pluginJson, _ := r.Get(fmt.Sprintf("github:jenkinsci:%s:%s", pluginName, pluginVersion)).Bytes()
	var releaseNote types.GitHubReleaseNote
	err := json.Unmarshal(pluginJson, &releaseNote)

	if err != nil {
		log.Println(err)
		// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
		return types.GitHubReleaseNote{}, errors.New("failed to unmarshal ReleaseNote")
	}
	return releaseNote, nil
}

func (r *Redis) GetPluginVersions(pluginName string) ([]byte, error) {

	pluginVersionsJson, err := r.Get(fmt.Sprintf("github:jenkinsci:%s:versions", pluginName)).Bytes()
	if err != nil {
		return []byte{}, errors.New("error in getPlugins " + pluginName)
	}
	if err != nil {
		// plugin doesnt exist
		fmt.Println("versions file is not exist")
		fmt.Println(err)
		// GetGitHubReleases(plugin.Name, redisclient)
		// plugin = redisclient.Get()
		return nil, err
	}
	return pluginVersionsJson, nil
}

func (r *Redis) SaveGithubStats(gh github.GitHubStats) error {
	jsonData, err := json.Marshal(gh)
	if err != nil {
		log.Println(err)
		fmt.Println("Failed to marshal GitHubStats")
	}
	// set lastUpdated file for repo
	err = r.Set("github:stats", jsonData)
	if err != nil {
		log.Println(err)
		return errors.New("set github:stats failed")
	}

	return nil
}
