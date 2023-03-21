package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

type projectsPage struct {
	Title      string
	Projects   []types.JenkinsPlugin
	ServerName string
}

func (h *RedisHandler) projectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/projects" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	var projects []types.JenkinsPlugin

	allProjects := h.Redis.GetAllPluginsFromServers()

	for _, projectName := range allProjects {
		projects = append(projects, types.JenkinsPlugin{
			Name:         projectName,
			IsDownloaded: h.Redis.IsProjectDownloaded(projectName),
			Error:        h.Redis.GetProjectError(projectName),
			LastUpdated:  h.Redis.GetLastUpdatedTime(projectName),
		})
	}

	pp := projectsPage{
		Title:    "Projects",
		Projects: projects,
	}
	// sp := ServersPage{
	// 	Title:   "Servers",
	// 	Servers: h.Redis.GetJenkinsServers(),
	// }
	h.Data = pp

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("web/templates/projects.html"))
	err := tmpl.Execute(w, h.Data)
	if err != nil {
		log.Println(err)
	}
}

func ProjectsHandler(redisHandler RedisHandler) {
	http.HandleFunc("/projects", redisHandler.projectsHandler)
}
