package redisStorage

import (
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/storage"
	"github.com/sirupsen/logrus"
)

type RedisStorage struct {
	DB storage.Storage
}

func (r *RedisStorage) AddDebugData() {
	logrus.Infoln("Append debug data to redis")
	// r.AddJenkinsServer("jenkins-two", "2.3233.1")
	// r.AddJenkinsServer("jenkins-one", "2.3233.2")

	// plugins := `mina-sshd-api-common:2.9.1-44.v476733c11f82`
	plugins := `
	ant:475.vf34069fef73c
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
