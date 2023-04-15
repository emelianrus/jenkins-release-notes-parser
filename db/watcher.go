package db

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

var WATCHER_LIST_PATH = "watcher:data"

func (r *Redis) GetWatcherData() (map[string]string, error) {
	watcherList, err := r.Get(WATCHER_LIST_PATH).Bytes()
	if err != nil {
		logrus.Errorln("can not get watcher list")
		logrus.Errorln(err)
	}
	var result map[string]string
	err = json.Unmarshal(watcherList, &result)
	if err != nil {
		logrus.Errorln("can not unmarshal watcherList")
		logrus.Errorln(err)
	}

	return result, nil
}

func (r *Redis) SetWatcherList(content map[string]string) error {
	jsonBody, err := json.Marshal(content)
	if err != nil {
		logrus.Errorln("failed to marshal body")
		logrus.Errorln(err)
	}

	err = r.Set(WATCHER_LIST_PATH, jsonBody)
	if err != nil {
		logrus.Errorln(err)
	}

	return nil
}