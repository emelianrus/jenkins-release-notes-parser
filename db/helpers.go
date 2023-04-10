package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/github"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
	"github.com/sirupsen/logrus"
)

func (r *Redis) GetJenkinsServers() []types.JenkinsServer {
	var servers []types.JenkinsServer

	keys, _ := r.Keys("servers:*")
	for _, path := range keys {
		re := strings.Split(path, ":")
		if len(re) == 2 {

			serverJson, _ := r.Get(path).Bytes()
			var jenkinsServer types.JenkinsServer
			err := json.Unmarshal(serverJson, &jenkinsServer)
			if err != nil {
				logrus.Errorln("can not unmarshal jenkins server")
			}

			projects, _ := r.GetJenkinsProjects(jenkinsServer.Name)
			servers = append(servers, types.JenkinsServer{
				Name:    jenkinsServer.Name,
				Core:    jenkinsServer.Core,
				Plugins: projects,
			})
		}
	}

	return servers
}

func (r *Redis) AddJenkinsServer(serverName string, coreVersion string) {
	// TODO: check if exist/replace
	// js := JenkinsServer{
	// 	Name: "jenkins-one",
	// 	Core: "2.3233.2",
	// 	Plugins: []ServerPlugin{
	// 		{
	// 			Name:    "plugin-installation-manager-tool",
	// 			Version: "2.10.0",
	// 		},
	// 		{
	// 			Name:    "okhttp-api-plugin",
	// 			Version: "4.9.3-108.v0feda04578cf",
	// 		},
	// 	},
	// }

	jsonData, err := json.Marshal(types.JenkinsServer{
		Name: serverName,
		Core: coreVersion,
	})

	if err != nil {
		logrus.Errorln(err)
		return
	}
	// write jenkins server json
	err = r.Set(fmt.Sprintf("servers:%s", serverName), jsonData)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	fmt.Printf("Added %s:%s\n", serverName, coreVersion)
}

func (r *Redis) DeleteJenkinsServer(serverName string) {
	path := fmt.Sprintf("servers:%s", serverName)
	logrus.Infof("Removing jenkins server %s\n", path)
	r.Del(path)
}

func (r *Redis) GetAllProjectsFromServers() []string {
	var result []string
	projectKeys, _ := r.Keys("servers:*:plugins:*")

	for _, path := range projectKeys {
		splitted := strings.Split(path, ":")
		result = append(result, splitted[len(splitted)-1])
	}
	result = utils.RemoveDuplicates(result)
	return result
}

func (r *Redis) ChangeJenkinServerPluginVersion(serverName string, pluginName string, newVersion string) error {
	logrus.Infof("Change plugin version to server name: %s plugin name %s new verion %s\n", serverName, pluginName, newVersion)
	jsonData, _ := json.Marshal(newVersion)

	err := r.Set(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName), jsonData)
	if err != nil {
		logrus.Errorln("can not append to redis new plugin")
	}
	logrus.Infof("version for plugin %s was changed and saved to redis\n", pluginName)

	return nil
}

func (r *Redis) GetJenkinsProjects(jenkinsServer string) ([]types.Project, error) {
	var jenkinsPlugins []types.Project
	projectKeys, _ := r.Keys(fmt.Sprintf("servers:%s:plugins:*", jenkinsServer))

	for _, projectKey := range projectKeys {

		// get plugin name from path (is there is a better way? )
		path := strings.Split(projectKey, ":")
		projectName := path[len(path)-1]
		pluginVersionByte, _ := r.Get(projectKey).Bytes()

		var version string
		err := json.Unmarshal(pluginVersionByte, &version)
		if err != nil {
			logrus.Errorln("can not unmarshal version")
		}
		projectError := r.GetProjectError(projectName)

		jenkinsPlugins = append(jenkinsPlugins, types.Project{
			Name:         projectName,
			Version:      version,
			Error:        projectError,
			IsDownloaded: r.IsProjectDownloaded(projectName),
		})
	}

	return jenkinsPlugins, nil
}

func (r *Redis) AddJenkinsServerPlugin(serverName string, plugin types.Project) error {
	_, err := r.Get(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name)).Bytes()
	if err != nil {
		logrus.Infoln(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name) + " not found in redis, adding...")

		jsonData, _ := json.Marshal(plugin.Version)

		err := r.Set(fmt.Sprintf("servers:%s:plugins:%s", serverName, plugin.Name), jsonData)
		if err != nil {
			logrus.Errorln("can not append to redis new plugin")
		}
		logrus.Infof("appended plugin %s to redis\n", plugin.Name)
	}
	return nil
}

func (r *Redis) RemoveJenkinsServerPlugin(serverName string, pluginName string) {
	logrus.Infof("removing key from redis %s\n", pluginName)
	r.Del(fmt.Sprintf("servers:%s:plugins:%s", serverName, pluginName))
}

func (r *Redis) SetLastUpdatedTime(pluginName string, value string) error {

	jsonData, _ := json.Marshal(value)
	err := r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "lastUpdated"),
		jsonData)
	if err != nil {
		logrus.Errorln("SetLastUpdatedTime error")
		logrus.Errorln(err)
		return nil
	}
	return nil
}

func (r *Redis) GetLastUpdatedTime(projectName string) string {

	serverJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "lastUpdated")).Bytes()
	return string(serverJson)
}

func (r *Redis) SetProjectError(projectName string, value string) error {
	jsonData, _ := json.Marshal(value)
	err := r.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "error"),
		jsonData)
	if err != nil {
		logrus.Errorln("SetProjectError error:")
		logrus.Errorln(err)
		return nil
	}
	return nil
}

func (r *Redis) GetProjectError(projectName string) string {
	serverJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "error")).Bytes()
	return string(serverJson)
}

func (r *Redis) IsProjectDownloaded(projectName string) bool {
	_, err := r.Get(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "versions")).Bytes()
	if err == nil {
		return true
	} else {
		return false
	}
}

// func (r *Redis) SetPluginWithVersion(pluginName string, pluginVersion string, releaseNote types.GitHubReleaseNote) error {
// 	if pluginName == "" || pluginVersion == "" {
// 		return fmt.Errorf("can not set to DB, name: %s or Version: %s is empty", pluginName, pluginVersion)
// 	}
// 	key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, pluginVersion)
// 	// 0 time.Hour
// 	jsonData, err := json.Marshal(releaseNote)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	err = r.Set(key, jsonData)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	return nil
// }

func (r *Redis) GetAllProject() ([]types.Project, error) {
	var projects []types.Project
	// TODO: this should be configurable
	repoOwner := "jenkinsci"

	// get list of all projects
	var projectsTmp []string
	projectsKeys, _ := r.Keys(fmt.Sprintf("github:%s:*", repoOwner))
	for _, key := range projectsKeys {
		splitted := strings.Split(key, ":")
		projectsTmp = append(projectsTmp, splitted[2])
	}
	projectsTmp = utils.RemoveDuplicates(projectsTmp)

	// gather project data
	for _, key := range projectsTmp {

		projects = append(projects, types.Project{
			Name:  key,
			Owner: repoOwner,
			// TODO: should gather all fields
			Error:        "some error",
			IsDownloaded: true, // r.IsProjectDownloaded(key)
			LastUpdated:  "Tue 15 2022",
		})

	}

	return projects, nil
}

func (r *Redis) GetProject(projectOwner string, projectName string) ([]types.GitHubReleaseNote, error) {

	releaseNotes := []types.GitHubReleaseNote{}

	// get all versions for specific project
	projectVersionsJson, err := r.Get(fmt.Sprintf("github:%s:%s:versions", projectOwner, projectName)).Bytes()
	if err != nil {
		return releaseNotes, errors.New("error in getPlugins " + projectName)
	}
	// convert json versions to []string
	var versions []string
	err = json.Unmarshal(projectVersionsJson, &versions)
	if err != nil {
		logrus.Errorln(err)
	}

	for _, version := range versions {
		// get release notes of specific release
		pluginJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, version)).Bytes()
		var releaseNote types.GitHubReleaseNote
		err := json.Unmarshal(pluginJson, &releaseNote)
		if err != nil {
			logrus.Errorln(err)
			// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
			return releaseNotes, errors.New("failed to unmarshal ReleaseNote")
		}
		// append release version note to list
		releaseNote.Body = string(template.HTML(
			utils.ReplaceGitHubLinks(
				utils.ConvertMarkDownToHtml(releaseNote.Body))))
		releaseNotes = append(releaseNotes, releaseNote)

	}

	return releaseNotes, nil
}

func (r *Redis) GetPluginWithVersion(projectName string, projectVersion string) (types.GitHubReleaseNote, error) {
	pluginJson, _ := r.Get(fmt.Sprintf("github:jenkinsci:%s:%s", projectName, projectVersion)).Bytes()
	var releaseNote types.GitHubReleaseNote
	err := json.Unmarshal(pluginJson, &releaseNote)

	if err != nil {
		logrus.Errorln(err)
		// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
		return types.GitHubReleaseNote{}, errors.New("failed to unmarshal ReleaseNote")
	}
	return releaseNote, nil
}

func (r *Redis) GetPluginVersions(projectName string) ([]byte, error) {

	projectVersionsJson, err := r.Get(fmt.Sprintf("github:jenkinsci:%s:versions", projectName)).Bytes()
	if err != nil {
		return []byte{}, errors.New("error in getPlugins " + projectName)
	}
	if err != nil {
		// plugin doesnt exist
		logrus.Errorln("versions file is not exist")
		logrus.Errorln(err)
		return nil, err
	}
	return projectVersionsJson, nil
}

func (r *Redis) SaveGithubStats(gh github.GitHubStats) error {
	jsonData, err := json.Marshal(gh)
	if err != nil {
		logrus.Errorln(err)
		logrus.Errorln("Failed to marshal GitHubStats")
	}
	// set lastUpdated file for repo
	err = r.Set("github:stats", jsonData)
	if err != nil {
		logrus.Errorln(err)
		return errors.New("set github:stats failed")
	}

	return nil
}
