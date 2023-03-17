package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

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

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 page not found")
	}
}

func (h *ServersHandler) serversHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

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

func handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "js/servers.js")
}

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

	http.HandleFunc("/release-notes", releaseNotesHandler.releaseNotesHandler)
	http.HandleFunc("/", serversHandler.serversHandler)
	http.HandleFunc("/js/", handleJS)

	http.HandleFunc("/delete-plugin", crudHandler.deleteJenkinsPlugin)

	// TODO: should be list
	// [{"name":"jenkinsServerName","value":"jenkins-two"},{"name":"pluginName","value":"1"},{"name":"pluginVersion","value":"2"},{"name":"pluginName","value":"3"},{"name":"pluginVersion","value":"4"}]
	// http.HandleFunc("/add-new-plugin", crudHandler.deleteJenkinsPlugin)
	// TODO:
	// http.HandleFunc("/change-plugin-version", crudHandler.deleteJenkinsPlugin)
	// TODO:
	// http.HandleFunc("/delete-server", crudHandler.deleteJenkinsPlugin)
	// TODO:
	// http.HandleFunc("/add-server", crudHandler.deleteJenkinsPlugin)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
