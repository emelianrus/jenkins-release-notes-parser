package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type RedisHandler struct {
	Redis *Redis
	Data  interface{}
}

// Handler for 404 page
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 page not found")
	}
}

type ReleaseNotesPage struct {
	Title      string
	Products   []Product
	ServerName string
}

func (h *RedisHandler) releaseNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve releases from Redis
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// GET methods
	query := r.URL.Query()
	pickedJenkinsName := query.Get("jenkins")

	jenkinsServers := h.Redis.getJenkinsServers()

	// iterate to get correct jenkins and plugins
	for _, jenkinsServer := range jenkinsServers {

		if jenkinsServer.Name == pickedJenkinsName {
			data, _ := getReleaseNotesPageData(h.Redis, jenkinsServer)
			releaseNotesData := ReleaseNotesPage{
				Title:      "Plugin manager",
				ServerName: jenkinsServer.Name,
				Products:   data,
			}
			h.Data = releaseNotesData
			break
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))
	err := tmpl.Execute(w, h.Data)
	if err != nil {
		log.Println(err)
	}
}

type ServersPage struct {
	Title   string
	Servers []JenkinsServer
}

func (h *RedisHandler) serversHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	sp := ServersPage{
		Title:   "Servers",
		Servers: h.Redis.getJenkinsServers(),
	}
	h.Data = sp

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/servers.html"))
	err := tmpl.Execute(w, h.Data)
	if err != nil {
		log.Println(err)
	}
}

type deleteJenkinsPluginPayload struct {
	JenkinsName string
	PluginName  string
}

// POST function to delete jenkins plugin from servers page
func (h *RedisHandler) deleteJenkinsPlugin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Read the request body
		// TODO: replace ioutil with io.Copy
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Do something with the request body, such as decoding it into a struct
		var data deleteJenkinsPluginPayload
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		h.Redis.removeJenkinsServerPlugin(data.JenkinsName, data.PluginName)

		// Do something with the data, such as processing it or saving it to a database

		// Send a response back to the client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request processed successfully"))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

type addJenkinsPluginPayload struct {
	JenkinsName string              `json:"jenkinsName"`
	Plugins     []map[string]string `json:"plugins"`
}

// POST add new jenkins plugin to jenkins server
func (h *RedisHandler) addJenkinsPlugin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var server addJenkinsPluginPayload
	err := decoder.Decode(&server)
	if err != nil {
		panic(err)
	}
	// TODO: Add check null field

	for _, plugin := range server.Plugins {

		for name, version := range plugin {
			plugin := JenkinsPlugin{
				Name:    fmt.Sprintf("%v", name),
				Version: fmt.Sprintf("%v", version),
			}
			fmt.Println(plugin)
			h.Redis.addJenkinsServerPlugin(server.JenkinsName, plugin)
		}
	}

}

type addJenkinsServerPayload struct {
	JenkinsName string
	CoreVersion string
}

// POST add jenkins server to DB
func (h *RedisHandler) addJenkinsServer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var server addJenkinsServerPayload
	err := decoder.Decode(&server)
	if err != nil {
		panic(err)
	}

	h.Redis.addJenkinsServer(server.JenkinsName, server.CoreVersion)
}

// POST delete jenkins server from DB
func (h *RedisHandler) deleteJenkinsServer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var serverName string
	err := decoder.Decode(&serverName)
	if err != nil {
		panic(err)
	}
	// w.WriteHeader(http.StatusBadRequest)
	h.Redis.deleteJenkinsServer(serverName)
}

type changePluginVersionPayload struct {
	JenkinsName      string
	PluginName       string
	NewPluginVersion string
}

func (h *RedisHandler) changePluginVersion(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newPlugin changePluginVersionPayload
	err := decoder.Decode(&newPlugin)
	if err != nil {
		panic(err)
	}
	h.Redis.changeJenkinServerPluginVersion(newPlugin.JenkinsName, newPlugin.PluginName, newPlugin.NewPluginVersion)
	// w.WriteHeader(http.StatusBadRequest)
}

// Load js from separate file
// TODO: Is there another way? i've found a lot of issues in http lib related to it
func handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "js/servers.js")
}

func StartWeb(redisclient *Redis) {
	redisHandler := RedisHandler{
		Redis: redisclient,
	}

	log.Println("Starting server")

	http.HandleFunc("/", redisHandler.serversHandler)
	http.HandleFunc("/release-notes", redisHandler.releaseNotesHandler)
	http.HandleFunc("/js/", handleJS)

	http.HandleFunc("/delete-plugin", redisHandler.deleteJenkinsPlugin)

	http.HandleFunc("/add-new-plugin", redisHandler.addJenkinsPlugin)
	http.HandleFunc("/change-plugin-version", redisHandler.changePluginVersion)
	http.HandleFunc("/delete-server", redisHandler.deleteJenkinsServer)
	http.HandleFunc("/add-server", redisHandler.addJenkinsServer)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
