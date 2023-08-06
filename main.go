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
	err := godotenv.Load()
	if err != nil {
		logrus.Errorln("Error loading .env file")
	}
	// Initial configuration for logger
	logLevel := os.Getenv("RN_DEBUG")
	if logLevel != "" {
		lvl, _ := logrus.ParseLevel(logLevel)
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if runtime.GOOS == "windows" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}
	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)
}

func Start() {

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
	err := router.Run(":8080")
	if err != nil {
		logrus.Errorln("Failed to create gin server")
		logrus.Errorln(err)
	}
}

func main() {
	Start()
}
