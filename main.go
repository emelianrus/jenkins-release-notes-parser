package main

import (
	"os"
	"runtime"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/routes"
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

	if redisclient.Status() != nil {
		logrus.Errorln("failed to connect to redis")
	} else {
		// TODO: remove used during development
		redisclient.AddDebugData()
	}

	// githubClient := github.NewGitHubClient()

	// TODO: should be update plugin function executed once per day
	// go worker.StartQueue(redisclient, githubClient)

	// GIN
	router := routes.SetupRouter()
	router.Run(":8080")
	// WEB
	// web.StartWeb(redisclient, &githubClient)
}
