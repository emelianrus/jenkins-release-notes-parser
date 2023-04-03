package main

import (
	"os"
	"runtime"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/github"
	"github.com/emelianrus/jenkins-release-notes-parser/web"
	"github.com/emelianrus/jenkins-release-notes-parser/worker"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if runtime.GOOS == "windows" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}
	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)
}

func main() {
	redisclient := db.NewRedisClient()
	githubClient := github.NewGitHubClient()
	// TODO: remove used during development
	redisclient.AddDebugData()

	// TODO: should be update plugin function executed once per day
	go worker.StartQueue(redisclient, githubClient)

	web.StartWeb(redisclient, &githubClient)
}
