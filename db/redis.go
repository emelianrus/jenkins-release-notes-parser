package db

import (
	"fmt"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedisClient() *Redis {
	fmt.Println("Creating redis connection")
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		// Password: "",
		DB: 0,
	})
	return &Redis{client: client}
}

func (r *Redis) AddDebugData() {
	r.AddJenkinsServer("jenkins-two", "2.3233.1")
	r.AddJenkinsServer("jenkins-one", "2.3233.2")

	r.AddJenkinsServerPlugin("jenkins-one", types.JenkinsPlugin{
		Name:    "plugin-installation-manager-tool",
		Version: "2.10.0",
	})
	r.AddJenkinsServerPlugin("jenkins-two", types.JenkinsPlugin{
		Name:    "plugin-installation-manager-tool",
		Version: "2.10.0",
	})
	r.AddJenkinsServerPlugin("jenkins-one", types.JenkinsPlugin{
		Name:    "okhttp-api-plugin",
		Version: "4.9.3-108.v0feda04578cf",
	})
}

func (r *Redis) Get(key string) *redis.StringCmd {
	return r.client.Get(key)
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}
