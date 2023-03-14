package main

import (
	"html/template"
	"log"
	"net/http"
)

// func constractPageData() {
// 	// getJenkinsInstance
// 	// getListOfPlugins
// 	// checkIsPluginInCache
// 	// createStructure
// }

type PluginPage struct {
	Title      string
	Products   []Product
	ServerName string
}

type PluginHandler struct {
	Data PluginPage
}

func (p *PluginHandler) pluginsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve releases from Redis
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))
	err := tmpl.Execute(w, p.Data)
	if err != nil {
		log.Println(err)
	}
}

type ServersPage struct {
	Title   string
	Servers []JenkinsServer
}

type ServersHandler struct {
	Data ServersPage
}

func (h *ServersHandler) serversHandler(w http.ResponseWriter, r *http.Request) {

	// plPage := ServersPage{
	// 	Title: "Plugin manager",
	// 	// ServerName: "jenkins-one",
	// 	// Products:   nil,
	// 	Servers: []JenkinsServer{
	// 		{
	// 			Name: "ser1",
	// 			Plugins: []JenkinsPlugin{
	// 				{
	// 					Name:    "pl1",
	// 					Version: "v1",
	// 				},
	// 				{
	// 					Name:    "pl2",
	// 					Version: "v33",
	// 				},
	// 			},
	// 		},
	// 		{
	// 			Name: "ser2",
	// 		},
	// 	},
	// }

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/test.html"))
	err := tmpl.Execute(w, h.Data)
	if err != nil {
		log.Println(err)
	}
}

func StartWeb(redisclient *Redis) {
	jenkinsServers := redisclient.getJenkinsServers()

	pluginHandler := PluginHandler{
		Data: getPluginsForPageData(redisclient, jenkinsServers[1]),
	}

	sp := ServersPage{
		Title:   "Servers",
		Servers: redisclient.getJenkinsServers(),
	}
	serversHandler := ServersHandler{
		Data: sp,
	}

	log.Println("Starting server")
	// data := getPlugins(redisclient)
	http.HandleFunc("/release-notes", pluginHandler.pluginsHandler)
	http.HandleFunc("/servers", serversHandler.serversHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
