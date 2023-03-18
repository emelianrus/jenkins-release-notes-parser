package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
)

// DB part end ^

func (r *Redis) GetJenkinsServers() []types.JenkinsServer {
	var servers []types.JenkinsServer

	keys, _ := r.client.Keys("servers:*").Result()
	for _, path := range keys {
		re := strings.Split(path, ":")
		if len(re) == 2 {

			serverJson, _ := r.client.Get(path).Bytes()
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
	r.client.Del(path)
}

func (r *Redis) GetAllPluginsFromServers() []string {
	var result []string
	pluginKeys, _ := r.client.Keys("servers:*:plugins:*").Result()

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
	fmt.Println("appended plugin to redis")

	return nil
}

func (r *Redis) GetJenkinsPlugins(jenkinsServer string) ([]types.JenkinsPlugin, error) {
	var jenkinsPlugins []types.JenkinsPlugin
	pluginKeys, _ := r.client.Keys(fmt.Sprintf("servers:%s:plugins:*", jenkinsServer)).Result()

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

		jenkinsPlugins = append(jenkinsPlugins, types.JenkinsPlugin{
			Name:    pluginName,
			Version: version,
		})
	}

	return jenkinsPlugins, nil
}

func (r *Redis) AddJenkinsServerPlugin(serverName string, plugin types.JenkinsPlugin) error {
	_, err := r.Get(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name)).Bytes()
	if err != nil {
		fmt.Println(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name) + " not found all good")

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
	// _, err := r.Get(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName)).Bytes()
	fmt.Printf("removing key from redis %s\n", pluginName)
	r.client.Del(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName))
}

// func (r *Redis) GetPlugin(key string) ([]byte, error) {
// 	jsonData, err := r.Get(key).Bytes()
// 	if err != nil {
// 		return []byte{}, errors.New("error in getPlugins " + key)
// 	}
// 	return jsonData, err
// }

func (r *Redis) SetLastUpdatedTime(pluginName string, value string) error {

	err := r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "lastUpdated"),
		value)
	if err != nil {
		log.Println("SetLastUpdatedTime error")
		log.Println(err)
		return nil
	}
	return nil
}

func (r *Redis) GetLastUpdatedTime(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) GetVersions(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) SetVersions(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) SetPluginWithVersion(pluginName string, pluginVersion string, releaseNote types.GitHubReleaseNote) error {
	key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, pluginVersion)
	// 0 time.Hour
	jsonData, err := json.Marshal(releaseNote)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = r.Set(key, jsonData)
	if err != nil {
		log.Println(err)
		return nil
	}
	return nil
}

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

func (r *Redis) SaveReleaseNotesToDB(releases []types.GitHubReleaseNote, pluginName string) error {

	currentTime := time.Now()
	formattedTime := currentTime.Format("02 January 2006 15:04")
	err := r.SetLastUpdatedTime(pluginName, formattedTime)
	if err != nil {
		// fmt.Println(err)
		// fmt.Println("Can not set updated time")
		return fmt.Errorf("error setting lastUpdate time in get github release: %s", err)
	}
	var versions []string

	for _, release := range releases {
		versions = append(versions, release.Name)
		key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, release.Name)
		// 0 time.Hour
		jsonData, err := json.Marshal(release)
		if err != nil {
			// log.Println(err)
			return fmt.Errorf("error Marshal release: %s", err)
		}
		err = r.Set(key, jsonData)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error setting release: %s", err)
		}
	}

	jsonVersions, _ := json.Marshal(versions)
	err = r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "versions"),
		jsonVersions)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}
	return nil
}
