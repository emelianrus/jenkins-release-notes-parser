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

func StartWeb(redisclient *Redis) {

	pluginHandler := PluginHandler{
		Data: getPlugins(redisclient),
	}

	log.Println("Starting server")
	// data := getPlugins(redisclient)
	http.HandleFunc("/", pluginHandler.pluginsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
