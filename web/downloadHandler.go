package web

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
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
		logrus.Errorln("Failed to get releases from github")
	}
	err = h.Redis.SaveReleaseNotesToDB(releases, projectName)
	if err != nil {
		logrus.Errorln(err)
		logrus.Errorln("Failed to save release notes to db")
	}

}

func DownloadHandler(redisHandler CommonHandler) {
	http.HandleFunc("/download", redisHandler.downloadHeandler)
}
