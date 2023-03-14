package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedisClient() *Redis {
	fmt.Println("Creating redis connection")
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		// Password: "",
		DB: 0,
	})
	return &Redis{client: client}
}

func (r *Redis) Get(key string) *redis.StringCmd {
	return r.client.Get(key)
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

// DB part end ^

type JenkinsPlugin struct {
	Name    string
	Version string
}

type JenkinsServer struct {
	Name    string
	Core    string
	Plugins []JenkinsPlugin
}

func (r *Redis) getJenkinsServers() []JenkinsServer {
	var servers []JenkinsServer

	keys, _ := r.client.Keys("servers:*").Result()
	for _, path := range keys {
		re := strings.Split(path, ":")
		if len(re) == 2 {

			serverJson, _ := r.client.Get(path).Bytes()
			var jenkinsServer JenkinsServer
			err := json.Unmarshal(serverJson, &jenkinsServer)
			if err != nil {
				fmt.Println("can not unmarshal jenkins server")
			}

			plugins, _ := r.getJenkinsPlugins(jenkinsServer.Name)
			servers = append(servers, JenkinsServer{
				Name:    jenkinsServer.Name,
				Core:    jenkinsServer.Core,
				Plugins: plugins,
			})
		}
	}

	return servers
}

func (r *Redis) addJenkinsServer(serverName string, coreVersion string) {
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

	jsonData, err := json.Marshal(JenkinsServer{
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
}

func (r *Redis) changeJenkinServerPluginVersion(serverName string, pluginName string, newVersion string) error {
	jsonData, _ := json.Marshal(newVersion)

	err := r.Set(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName), jsonData)
	if err != nil {
		fmt.Println("can not append to redis new plugin")
	}
	fmt.Println("appended plugin to redis")

	return nil
}

func (r *Redis) getJenkinsPlugins(jenkinsServer string) ([]JenkinsPlugin, error) {
	var jenkinsPlugins []JenkinsPlugin
	pluginKeys, _ := r.client.Keys(fmt.Sprintf("servers:%s:plugins:*", jenkinsServer)).Result()

	for _, pluginKey := range pluginKeys {

		// get plugin name from path (is there is a better way? )
		path := strings.Split(pluginKey, ":")
		pluginName := path[len(path)-1]
		pluginVersionByte, _ := r.Get(pluginKey).Bytes()

		pluginVersion := string(pluginVersionByte)

		jenkinsPlugins = append(jenkinsPlugins, JenkinsPlugin{
			Name:    pluginName,
			Version: pluginVersion,
		})
	}

	return jenkinsPlugins, nil
}

func (r *Redis) addJenkinsServerPlugin(serverName string, plugin JenkinsPlugin) error {
	_, err := r.Get(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name)).Bytes()
	if err != nil {
		fmt.Println(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name) + " not found all good")

		jsonData, _ := json.Marshal(plugin.Version)

		err := r.Set(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name), jsonData)
		if err != nil {
			fmt.Println("can not append to redis new plugin")
		}
		fmt.Println("appended plugin to redis")
	}
	return nil
}

func (r *Redis) removeJenkinsServerPlugin(serverName string, pluginName string) {
	// _, err := r.Get(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName)).Bytes()
	fmt.Printf("removing key from redis %s\n", pluginName)
	r.client.Del(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName))
}

func (r *Redis) GetJenkinsServers() ([]byte, error) {
	return r.client.Get("servers:jenkins-one:plugins").Bytes()
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

func (r *Redis) SetPluginWithVersion(pluginName string, pluginVersion string, releaseNote GitHubReleaseNote) error {
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

func (r *Redis) GetPluginWithVersion(pluginName string, pluginVersion string) (GitHubReleaseNote, error) {
	pluginJson, _ := r.Get(fmt.Sprintf("github:jenkinsci:%s:%s", pluginName, pluginVersion)).Bytes()
	var releaseNote GitHubReleaseNote
	err := json.Unmarshal(pluginJson, &releaseNote)

	if err != nil {
		log.Println(err)
		// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
		return GitHubReleaseNote{}, errors.New("failed to unmarshal ReleaseNote")
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
