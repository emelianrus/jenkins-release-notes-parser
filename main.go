package main

import "github.com/emelianrus/jenkins-release-notes-parser/db"

func main() {
	redisclient := db.NewRedisClient()

	// TODO: remove used during development
	redisclient.AddDebugData()

	StartWeb(redisclient)
}
