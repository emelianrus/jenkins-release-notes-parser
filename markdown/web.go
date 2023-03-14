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

type ReleaseNotesPage struct {
	Title      string
	Products   []Product
	ServerName string
}

type ReleaseNotesHandler struct {
	Redis      *Redis
	AllServers []JenkinsServer
	Data       ReleaseNotesPage
}

func (p *ReleaseNotesHandler) pluginsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve releases from Redis
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	query := r.URL.Query()
	pickedJenkinsName := query.Get("jenkins")

	jenkinsServers := p.Redis.getJenkinsServers()
	for _, v := range jenkinsServers {

		if v.Name == pickedJenkinsName {
			p.Data = getPluginsForPageData(p.Redis, v)
		}
	}

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/servers.html"))
	err := tmpl.Execute(w, h.Data)
	if err != nil {
		log.Println(err)
	}
}

func StartWeb(redisclient *Redis) {
	// jenkinsServers := redisclient.getJenkinsServers()

	releaseNotesHandler := ReleaseNotesHandler{
		Redis: redisclient,
		// AllServers: jenkinsServers,
		// Data: getPluginsForPageData(redisclient, jenkinsServers[1]),
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
	http.HandleFunc("/release-notes", releaseNotesHandler.pluginsHandler)
	http.HandleFunc("/", serversHandler.serversHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
