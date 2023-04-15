package redisStorage

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (r *RedisStorage) SetLastUpdatedTime(pluginName string, value string) error {
	logrus.Debugf("update latestUpdated time to: %v", value)

	jsonData, _ := json.Marshal(value)
	err := r.DB.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "lastUpdated"),
		jsonData)
	if err != nil {
		logrus.Errorln("SetLastUpdatedTime error")
		logrus.Errorln(err)
		return nil
	}
	return nil
}

func (r *RedisStorage) GetLastUpdatedTime(projectOwner string, projectName string) string {

	serverJson, _ := r.DB.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, "lastUpdated"))
	return string(serverJson)
}

func (r *RedisStorage) SetProjectError(projectName string, value string) error {
	jsonData, _ := json.Marshal(value)
	err := r.DB.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", projectName, "error"),
		jsonData)
	if err != nil {
		logrus.Errorln("SetProjectError error:")
		logrus.Errorln(err)
		return nil
	}
	return nil
}

func (r *RedisStorage) GetProjectError(projectOwner string, projectName string) string {
	serverJson, _ := r.DB.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, "error"))
	return string(serverJson)
}

func (r *RedisStorage) IsProjectDownloaded(projectOwner string, projectName string) bool {
	_, err := r.DB.Get(fmt.Sprintf("github:%s:%s:%s", projectOwner, projectName, "versions"))
	if err == nil {
		return true
	} else {
		return false
	}
}
