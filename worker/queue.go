package worker

import (
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/sources"
	jenkins "github.com/emelianrus/jenkins-release-notes-parser/pkg/sources/jenkinsPluginSite"
	"github.com/emelianrus/jenkins-release-notes-parser/storage/redisStorage"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/sirupsen/logrus"
)

// var (
// 	serviceMutex sync.Mutex
// )

// can be executed by button from UI so we need to be sure running only one instance at once
func StartWorkerPluginSite(redisclient *redisStorage.RedisStorage, ps jenkins.PluginSite) {
	logrus.Infoln("StartQueue...")

	projects, _ := redisclient.GetPluginListData()
	for {
		for projectName := range projects {

			releaseNotes, err := sources.DownloadProjectReleaseNotes(&ps, projectName)
			if err != nil {
				releaseNotes = []types.ReleaseNote{}
				redisclient.SetProjectError(projectName, err.Error())
				logrus.Errorln("Downloading repo error:")
				logrus.Errorln(err)
			} else {
				latestVersionInCache, _ := redisclient.GetLatestVersion("jenkinsci", projectName)
				logrus.Infof("latestVersionInCache: %s releaseNotes: %s\n", latestVersionInCache, releaseNotes[0].Name)
			}

			redisclient.SaveReleaseNotesToDB(releaseNotes, projectName)
		}

		logrus.Infoln("StartWorkerPluginSite done. doing sleep for 24h")
		time.Sleep(time.Hour * 24)
	}

}

// func StartQueuePluginsSite(redisclient *db.Redis, gh github.GitHub, ps jenkins.PluginSite) {
// 	logrus.Infoln("StartQueuePluginsSite...")

// 	for {
// 		projects, err := redisclient.GetAllProjects()
// 		fmt.Println(projects)
// 		if err != nil {
// 			logrus.Errorln("can not get projects")
// 		}
// 		for _, project := range projects {
// 			// TODO: error api 404
// 			releaseNotes, err := sources.DownloadPlugin(&ps, project.Name)
// 			// releaseNotes, err := gh.Download(projectName)

// 			if err == nil {
// 				redisclient.SaveReleaseNotesToDB(releaseNotes, project.Name)
// 			} else {
// 				logrus.Errorln("Downloading repo error:")
// 				logrus.Errorln(err)
// 				redisclient.SaveReleaseNotesToDB([]types.ReleaseNote{}, project.Name)
// 				redisclient.SetProjectError(project.Name, err.Error())
// 			}

// 			// redisclient.SaveGithubStats(gh.GitHubStats)
// 			redisclient.SaveReleaseNotesToDB(releaseNotes, project.Name)
// 		}
// 		logrus.Infoln("sleep 3 hours")
// 		time.Sleep(time.Hour * 12)
// 	}
// }

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
