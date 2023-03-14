package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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

func (p *ReleaseNotesHandler) releaseNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve releases from Redis
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	query := r.URL.Query()
	pickedJenkinsName := query.Get("jenkins")

	jenkinsServers := p.Redis.getJenkinsServers()
	for _, v := range jenkinsServers {

		if v.Name == pickedJenkinsName {
			p.Data = getPluginsForPageData(p.Redis, v)
			fmt.Println(p.Data.Products)
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
	Redis *Redis
	Data  ServersPage
}

func (h *ServersHandler) serversHandler(w http.ResponseWriter, r *http.Request) {
	sp := ServersPage{
		Title:   "Servers",
		Servers: h.Redis.getJenkinsServers(),
	}
	h.Data = sp
	// serversHandler := ServersHandler{
	// 	Data: sp,
	// }
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/servers.html"))
	err := tmpl.Execute(w, h.Data)
	if err != nil {
		log.Println(err)
	}
}

type createDeleteHandler struct {
	redis *Redis
}

type deleteJenkinsPlugin struct {
	JenkinsName string
	PluginName  string
}

func (h *createDeleteHandler) deleteJenkinsPlugin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Do something with the request body, such as decoding it into a struct
		var data deleteJenkinsPlugin
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		h.redis.removeJenkinsServerPlugin(data.JenkinsName, data.PluginName)

		// Do something with the data, such as processing it or saving it to a database

		// Send a response back to the client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request processed successfully"))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (h *createDeleteHandler) addJenkinsPlugin(w http.ResponseWriter, r *http.Request) {}

func (h *createDeleteHandler) addJenkinsServer(w http.ResponseWriter, r *http.Request) {}

func (h *createDeleteHandler) deleteJenkinsServer(w http.ResponseWriter, r *http.Request) {}
func (h *createDeleteHandler) changePluginVersion(w http.ResponseWriter, r *http.Request) {}

func StartWeb(redisclient *Redis) {
	// jenkinsServers := redisclient.getJenkinsServers()

	releaseNotesHandler := ReleaseNotesHandler{
		Redis: redisclient,
		// AllServers: jenkinsServers,
		// Data: getPluginsForPageData(redisclient, jenkinsServers[1]),
	}

	// sp := ServersPage{
	// 	Title:   "Servers",
	// 	Servers: redisclient.getJenkinsServers(),
	// }
	serversHandler := ServersHandler{
		Redis: redisclient,
	}

	crudHandler := createDeleteHandler{
		redis: redisclient,
	}

	log.Println("Starting server")
	// data := getPlugins(redisclient)
	http.HandleFunc("/release-notes", releaseNotesHandler.releaseNotesHandler)
	http.HandleFunc("/", serversHandler.serversHandler)

	// TODO:
	http.HandleFunc("/delete-plugin", crudHandler.deleteJenkinsPlugin)
	// http.HandleFunc("/add-plugin", ssss)
	// http.HandleFunc("/add-jenkins-server", ssss)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
