package db

import (
	"strings"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	client *redis.Client
}

func NewRedisClient() *Redis {
	logrus.Infoln("Creating redis connection")
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		// Password: "",
		DB: 0,
	})
	return &Redis{client: client}
}

func (r *Redis) Status() error {
	return r.client.Ping().Err()
}

func (r *Redis) AddDebugData() {
	logrus.Infoln("Append debug data to redis")
	// r.AddJenkinsServer("jenkins-two", "2.3233.1")
	// r.AddJenkinsServer("jenkins-one", "2.3233.2")

	// plugins := `mina-sshd-api-common:2.9.1-44.v476733c11f82`
	plugins := `
	ant:481.v7b_09e538fcca
	antisamy-markup-formatter:159.v25b_c67cd35fb_`

	lines := strings.Split(plugins, "\n")

	m := make(map[string]string)
	for _, line := range lines {
		pair := strings.Split(line, ":")
		if len(pair) == 2 {
			m[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}

	r.SetWatcherList(m)

}

func (r *Redis) Get(key string) *redis.StringCmd {
	return r.client.Get(key)
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) Keys(key string) ([]string, error) {
	return r.client.Keys(key).Result()
}

func (r *Redis) Del(key string) *redis.IntCmd {
	return r.client.Del(key)
}
