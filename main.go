package main

import (
	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/web"
)

func main() {
	redisclient := db.NewRedisClient()

	// TODO: remove used during development
	redisclient.AddDebugData()

	// TODO: should be update plugin function executed once per day
	// go github.StartQueue(redisclient)

	web.StartWeb(redisclient)
}
