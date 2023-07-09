package main

import (
	"os"
	"runtime"

	"github.com/emelianrus/jenkins-release-notes-parser/routes"
	"github.com/emelianrus/jenkins-release-notes-parser/storage"
	"github.com/emelianrus/jenkins-release-notes-parser/storage/db"
	rs "github.com/emelianrus/jenkins-release-notes-parser/storage/redisStorage"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	// Initial configuration for logger
	logrus.SetLevel(logrus.DebugLevel)
	if runtime.GOOS == "windows" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}
	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)
}

func Start() {

	// read env vars
	err := godotenv.Load()
	if err != nil {
		logrus.Errorln("Error loading .env file")
	}

	redis := db.NewRedisClient()

	redisStorage := &rs.RedisStorage{
		DB: storage.SetStorage(redis),
	}

	// values for debug
	if redis.Status() != nil {
		logrus.Errorln("failed to connect to redis")
	} else {
		// TODO: remove used during development
		redisStorage.AddDebugData()
	}

	// githubClient := github.NewGitHubClient()

	// TODO: should be update plugin function executed once per day
	// go worker.StartWorkerPluginSite(redisStorage, jenkins.NewPluginSite())

	// GIN
	router := routes.SetupRouter(redisStorage)
	err = router.Run(":8080")
	if err != nil {
		logrus.Errorln("Failed to create gin server")
		logrus.Errorln(err)
	}
}

func main() {
	Start()
}
