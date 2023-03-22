package web

import (
	"html/template"
	"log"
	"net/http"
)

type ReleaseNotesPage struct {
	Title         string
	GitHubProject []GitHubProject
	ServerName    string
}

func (h *CommonHandler) releaseNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve releases from Redis
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// GET methods
	query := r.URL.Query()
	pickedJenkinsName := query.Get("jenkins")

	jenkinsServers := h.Redis.GetJenkinsServers()

	// iterate to get correct jenkins and plugins
	for _, jenkinsServer := range jenkinsServers {

		if jenkinsServer.Name == pickedJenkinsName {
			data, _ := getReleaseNotesPageData(h.Redis, jenkinsServer)
			releaseNotesData := ReleaseNotesPage{
				Title:         "Plugin manager",
				ServerName:    jenkinsServer.Name,
				GitHubProject: data,
			}
			h.Data = releaseNotesData
			break
		}
	}

	tmpl := template.Must(template.ParseFiles("web/templates/release-notes.html"))
	err := tmpl.Execute(w, h.Data)
	if err != nil {
		log.Println(err)
	}
}

func ReleaseNotesHandler(redisHandler CommonHandler) {
	http.HandleFunc("/release-notes", redisHandler.releaseNotesHandler)
}
