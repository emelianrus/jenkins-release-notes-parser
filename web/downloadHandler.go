package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/github"
)

// POST add jenkins server to DB
func (h *RedisHandler) downloadHeandler(w http.ResponseWriter, r *http.Request) {
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

	releases, err := github.Download(projectName)
	if err != nil {
		fmt.Println("Failed to get releases from github")
	}
	err = h.Redis.SaveReleaseNotesToDB(releases, projectName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to save release notes to db")
	}

}

func DownloadHandler(redisHandler RedisHandler) {
	http.HandleFunc("/download", redisHandler.downloadHeandler)
}
