package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

type RedisHandler struct {
	Redis *db.Redis
	Data  interface{}
}

// Handler for 404 page
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 page not found")
	}
}

type ServersPage struct {
	Title   string
	Servers []types.JenkinsServer
}

func (h *RedisHandler) serversHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	sp := ServersPage{
		Title:   "Servers",
		Servers: h.Redis.GetJenkinsServers(),
	}
	h.Data = sp

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("web/templates/servers.html"))
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
		h.Redis.RemoveJenkinsServerPlugin(data.JenkinsName, data.PluginName)

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

			// validation
			if name == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Println(w, "name is empty its not possible to handle")
				continue
			}

			plugin := types.JenkinsPlugin{
				Name:    fmt.Sprintf("%v", name),
				Version: fmt.Sprintf("%v", version),
			}

			// // pluginVersionsJson, err = r.GetPluginVersions(plugin.Name)
			// // if err != nil {
			// // 	fmt.Println(err)
			// // 	fmt.Println("2nd attempt to GetPluginVersions failed")
			// // 	// return web page with default values
			// // 	return []Product{}, errors.New("2nd attempt to GetPluginVersions failed")
			// // }

			// releases, err := github.Download(plugin.Name)
			// if err == nil {
			// 	h.Redis.SaveReleaseNotesToDB(releases, plugin.Name)
			// } else {
			// 	h.Redis.SaveReleaseNotesToDB([]types.GitHubReleaseNote{}, plugin.Name)
			// }

			h.Redis.AddJenkinsServerPlugin(server.JenkinsName, plugin)
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

	h.Redis.AddJenkinsServer(server.JenkinsName, server.CoreVersion)
}

// POST delete jenkins server from DB
func (h *RedisHandler) deleteJenkinsServer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	// could be only number, need to do validation
	var serverName interface{}
	err := decoder.Decode(&serverName)
	if err != nil {
		panic(err)
	}

	// w.WriteHeader(http.StatusBadRequest)
	h.Redis.DeleteJenkinsServer(fmt.Sprintf("%v", serverName))
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
	h.Redis.ChangeJenkinServerPluginVersion(newPlugin.JenkinsName, newPlugin.PluginName, newPlugin.NewPluginVersion)
	// w.WriteHeader(http.StatusBadRequest)
}

// Load js from separate file
// TODO: Is there another way? i've found a lot of issues in http lib related to it
func handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "web/js/servers.js")
}

func StartWeb(redisclient *db.Redis) {
	redisHandler := RedisHandler{
		Redis: redisclient,
	}

	log.Println("Starting server")

	// Pages
	http.HandleFunc("/", redisHandler.serversHandler)
	http.HandleFunc("/js/", handleJS)

	ReleaseNotesHandler(redisHandler)

	// POST handlers
	http.HandleFunc("/add-new-plugin", redisHandler.addJenkinsPlugin)
	http.HandleFunc("/delete-plugin", redisHandler.deleteJenkinsPlugin)

	http.HandleFunc("/change-plugin-version", redisHandler.changePluginVersion)

	http.HandleFunc("/add-server", redisHandler.addJenkinsServer)
	http.HandleFunc("/delete-server", redisHandler.deleteJenkinsServer)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
