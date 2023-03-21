package github

import (
	"fmt"
	"sync"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
)

var (
	serviceMutex sync.Mutex
)

// used as go StartQueue()
// can be executed by button from UI so we need to be sure running only one instance at once
// func StartQueue(redisclient *db.Redis) {
// 	serviceMutex.Lock()
// 	defer serviceMutex.Unlock()

// 	for {
// 		for _, pluginName := range redisclient.GetAllPluginsFromServers() {
// 			// TODO: error api 404
// 			ghReleaseNotes, _ := Download(pluginName)
// 			redisclient.SaveReleaseNotesToDB(ghReleaseNotes, pluginName)
// 		}
// 		fmt.Println("sleep 3 hours")
// 		time.Sleep(time.Hour * 3)
// 	}
// }

func StartQueue(redisclient *db.Redis, plugins []string, infinite bool) {
	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	fmt.Printf("started queue infinite: %t \n", infinite)

	for {
		for _, pluginName := range redisclient.GetAllPluginsFromServers() {
			// TODO: error api 404
			ghReleaseNotes, _ := Download(pluginName)
			redisclient.SaveReleaseNotesToDB(ghReleaseNotes, pluginName)
		}

		if !infinite {
			break
		}

		fmt.Println("sleep 3 hours")
		time.Sleep(time.Hour * 3)
	}
}
