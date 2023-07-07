package redisStorage

import (
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/storage"
	"github.com/sirupsen/logrus"
)

type RedisStorage struct {
	DB storage.Storage
}

// TODO: temp function, this should be set via API
func (r *RedisStorage) AddDebugData() {
	logrus.Infoln("Append debug data to redis")

	plugins := `
	ant:475.vf34069fef73c
	antisamy-markup-formatter:159.v25b_c67cd35fb_
	kubernetes:3622.va_9dc5592b_10c
	kubernetes-client-api:5.4.2`

	lines := strings.Split(plugins, "\n")

	m := make(map[string]string)
	for _, line := range lines {
		pair := strings.Split(line, ":")
		if len(pair) == 2 {
			m[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}

	r.SetWatcherData(m)

}
