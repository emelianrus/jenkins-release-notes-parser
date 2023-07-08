package handlers

import (
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/pluginManager"
	"github.com/emelianrus/jenkins-release-notes-parser/storage/redisStorage"
)

// struct for handlers to use DB connection
type ProjectService struct {
	Redis         *redisStorage.RedisStorage
	PluginManager pluginManager.PluginManager
}
