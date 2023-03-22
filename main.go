package main

import (
	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/github"
	"github.com/emelianrus/jenkins-release-notes-parser/web"
	"github.com/emelianrus/jenkins-release-notes-parser/worker"
)

func main() {
	redisclient := db.NewRedisClient()
	githubClient := github.NewGitHubClient()
	// TODO: remove used during development
	redisclient.AddDebugData()

	// TODO: should be update plugin function executed once per day
	go worker.StartQueue(redisclient, githubClient)

	web.StartWeb(redisclient, &githubClient)
}
