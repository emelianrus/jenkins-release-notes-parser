package redisStorage

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

var PLUGIN_LIST_PATH = "plugin-list:data"

func (r *RedisStorage) GetPluginListData() (map[string]string, error) {
	pluginList, err := r.DB.Get(PLUGIN_LIST_PATH)
	if err != nil {
		logrus.Errorln("can not get plugin list")
		logrus.Errorln(err)
	}
	var result map[string]string
	err = json.Unmarshal(pluginList, &result)
	if err != nil {
		logrus.Errorln("can not unmarshal pluginList")
		logrus.Errorln(err)
	}

	return result, nil
}

func (r *RedisStorage) SetPluginListData(content map[string]string) error {
	jsonBody, err := json.Marshal(content)
	if err != nil {
		logrus.Errorln("failed to marshal body")
		logrus.Errorln(err)
	}

	err = r.DB.Set(PLUGIN_LIST_PATH, jsonBody)
	if err != nil {
		logrus.Errorln(err)
	}

	return nil
}
