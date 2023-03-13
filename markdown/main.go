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
	redisclient := InitDB()

	// init jenkins server
	js := JenkinsServer{
		Name: "jenkins-one",
		Core: "2.3233.2",
		Plugins: []ServerPlugin{
			{
				Name:    "plugin-installation-manager-tool",
				Version: "2.12.9",
			},
		},
	}

	jsonData, err := json.Marshal(js)
	if err != nil {
		log.Println(err)
		return
	}
	// write jenkins server json
	err = redisclient.Set("servers:jenkins-one:plugins", jsonData, 0).Err()
	if err != nil {
		log.Println(err)
		return
	}

	Parser(redisclient)
	StartWeb(redisclient)
}
