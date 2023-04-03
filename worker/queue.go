package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/github"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

var (
	serviceMutex sync.Mutex
)

// used as go StartQueue()
// can be executed by button from UI so we need to be sure running only one instance at once
func StartQueue(redisclient *db.Redis, github github.GitHub) {
	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	for {
		for _, pluginName := range redisclient.GetAllPluginsFromServers() {
			// TODO: error api 404
			ghReleaseNotes, err := github.Download(pluginName)

			if err == nil {
				redisclient.SaveReleaseNotesToDB(ghReleaseNotes, pluginName)
			} else {
				fmt.Println("Downloading repo error:")
				fmt.Println(err)
				redisclient.SaveReleaseNotesToDB([]types.GitHubReleaseNote{}, pluginName)
				redisclient.SetProjectError(pluginName, err.Error())
			}

			redisclient.SaveGithubStats(github.GitHubStats)
			redisclient.SaveReleaseNotesToDB(ghReleaseNotes, pluginName)
		}
		fmt.Println("sleep 3 hours")
		time.Sleep(time.Hour * 3)
	}
}

// func StartQueue(redisclient db.Redis, githubClient GitHub, plugins []string, infinite bool) {
// 	serviceMutex.Lock()
// 	defer serviceMutex.Unlock()

// 	fmt.Printf("started queue infinite: %t \n", infinite)

// 	for {
// 		for _, pluginName := range redisclient.GetAllPluginsFromServers() {
// 			// TODO: error api 404
// 			ghReleaseNotes, _ := githubClient.Download(pluginName)
// 			redisclient.SaveReleaseNotesToDB(ghReleaseNotes, pluginName)
// 		}

// 		if !infinite {
// 			break
// 		}

// 		fmt.Println("sleep 3 hours")
// 		time.Sleep(time.Hour * 3)
// 	}
// }
