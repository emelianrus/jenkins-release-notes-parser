package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

type projectsPage struct {
	Title      string
	Projects   []types.Project
	ServerName string
}

func (h *CommonHandler) projectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/projects" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	var projects []types.Project

	allProjects := h.Redis.GetAllProjectsFromServers()

	for _, projectName := range allProjects {
		projects = append(projects, types.Project{
			Name:         projectName,
			IsDownloaded: h.Redis.IsProjectDownloaded("jenkinsci", projectName),
			Error:        h.Redis.GetProjectError("jenkinsci", projectName),
			LastUpdated:  h.Redis.GetLastUpdatedTime("jenkinsci", projectName),
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

func ProjectsHandler(redisHandler CommonHandler) {
	http.HandleFunc("/projects", redisHandler.projectsHandler)
}
