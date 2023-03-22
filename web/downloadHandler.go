package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// POST add jenkins server to DB
func (h *CommonHandler) downloadHeandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/download" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var projectName string
	err := decoder.Decode(&projectName)
	if err != nil {
		panic(err)
	}

	releases, err := h.GitHub.Download(projectName)
	h.Redis.SaveGithubStats(h.GitHub.GitHubStats)
	if err != nil {
		fmt.Println("Failed to get releases from github")
	}
	err = h.Redis.SaveReleaseNotesToDB(releases, projectName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to save release notes to db")
	}

}

func DownloadHandler(redisHandler CommonHandler) {
	http.HandleFunc("/download", redisHandler.downloadHeandler)
}
