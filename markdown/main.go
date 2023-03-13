package main

import (
	"encoding/json"
	"log"
)

type ServerPlugin struct {
	Name    string
	Version string
}

type JenkinsServer struct {
	Name    string
	Core    string
	Plugins []ServerPlugin
}

func main() {
	redisclient := NewRedisClient()

	// init jenkins server
	js := JenkinsServer{
		Name: "jenkins-one",
		Core: "2.3233.2",
		Plugins: []ServerPlugin{
			{
				Name:    "plugin-installation-manager-tool",
				Version: "2.10.0",
			},
		},
	}

	jsonData, err := json.Marshal(js)
	if err != nil {
		log.Println(err)
		return
	}
	// write jenkins server json
	err = redisclient.Set("servers:jenkins-one:plugins", jsonData)
	if err != nil {
		log.Println(err)
		return
	}

	// getPlugins(redisclient)
	// // GetPluginFromGitHub(redisclient)
	StartWeb(redisclient)

	// plugin, err := redisclient.GetPlugin("github:jenkinsci:plugin-installation-manager-tool:versions")

	// fmt.Println(err)
	// fmt.Println(plugin)
}
