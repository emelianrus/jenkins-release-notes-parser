package db

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (r *Redis) SetLastUpdatedTime(pluginName string, value string) error {
	logrus.Debugf("update latestUpdated time to: %v", value)

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

func (r *Redis) GetLastUpdatedTime(projectOwner string, projectName string) string {

	serverJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, "lastUpdated")).Bytes()
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

func (r *Redis) GetProjectError(projectOwner string, projectName string) string {
	serverJson, _ := r.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, "error")).Bytes()
	return string(serverJson)
}

func (r *Redis) IsProjectDownloaded(projectOwner string, projectName string) bool {
	_, err := r.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, "versions")).Bytes()
	if err == nil {
		return true
	} else {
		return false
	}
}
