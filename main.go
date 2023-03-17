package main

import (
	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/web"
)

func main() {
	redisclient := db.NewRedisClient()

	// TODO: remove used during development
	redisclient.AddDebugData()

	web.StartWeb(redisclient)
}
