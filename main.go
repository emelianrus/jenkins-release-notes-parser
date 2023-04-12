package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/routes"
	"github.com/emelianrus/jenkins-release-notes-parser/sources"
	"github.com/emelianrus/jenkins-release-notes-parser/sources/github"
	jenkins "github.com/emelianrus/jenkins-release-notes-parser/sources/jenkinsPluginSite"
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

func Start() {
	redisclient := db.NewRedisClient()

	if redisclient.Status() != nil {
		logrus.Errorln("failed to connect to redis")
	} else {
		// TODO: remove used during development
		redisclient.AddDebugData()
	}

	// githubClient := github.NewGitHubClient()

	// githubClient := github.NewGitHubClient()
	// pluginSiteClient := jenkins.NewPluginSite()

	// TODO: should be update plugin function executed once per day
	// go worker.StartQueue(redisclient, githubClient, pluginSiteClient)

	// GIN
	router := routes.SetupRouter(redisclient)
	router.Run(":8080")
	// WEB
	// web.StartWeb(redisclient, &githubClient)
}

func Testing() {

	github := github.NewGitHubClient()
	pluginSite := jenkins.NewPluginSite()

	releases, err := sources.DownloadPlugin(&pluginSite, "ant")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range releases {
		fmt.Println(v.Name)
	}

	releases, err = sources.DownloadPlugin(&github, "ant")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range releases {
		fmt.Println(v.Name)
	}

}

func main() {
	Start()
	// Testing()

}
