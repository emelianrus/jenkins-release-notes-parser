package main

import (
	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/github"
	"github.com/emelianrus/jenkins-release-notes-parser/web"
)

func main() {
	redisclient := db.NewRedisClient()

	// TODO: remove used during development
	redisclient.AddDebugData()

	go github.StartQueue(redisclient)

	web.StartWeb(redisclient)
}
