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

func serversHandler(w http.ResponseWriter, r *http.Request) {

	// data:=
	plPage := PluginPage{
		Title:      "Plugin manager",
		ServerName: "jenkins-one",
		Products:   nil,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/servers.html"))
	err := tmpl.Execute(w, plPage)
	if err != nil {
		log.Println(err)
	}
}

func StartWeb(redisclient *Redis) {
	jenkinsServers := redisclient.getJenkinsServers()

	pluginHandler := PluginHandler{
		Data: getPluginsForPageData(redisclient, jenkinsServers[1]),
	}

	log.Println("Starting server")
	// data := getPlugins(redisclient)
	http.HandleFunc("/", pluginHandler.pluginsHandler)
	http.HandleFunc("/servers", serversHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
