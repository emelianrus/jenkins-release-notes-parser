package db

import (
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
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
	// r.AddJenkinsServer("jenkins-two", "2.3233.1")
	r.AddJenkinsServer("jenkins-one", "2.3233.2")

	plugins := `
	ace-editor:1.1
	ansicolor:1.0.2
	ant:481.v7b_09e538fcca
	antisamy-markup-formatter:2.7
	apache-httpcomponents-client-4-api:4.5.13-138.v4e7d9a_7b_a_e61
	jackson2-api:2.13.4-293.vee957901b_6fb
	javadoc:226.v71211feb_e7e9
	jdk-tool:55.v1b_32b_6ca_f9ca
	jenkins-design-language:1.25.8
	jira:3.8
	pipeline-graph-analysis:195.v5812d95a_a_2f9
	pipeline-maven:1205.vceea_7b_972817
	windows-slaves:1.8.1
	workflow-aggregator:590.v6a_d052e5a_a_b_5
	workflow-api:1200.v8005c684b_a_c6
	workflow-basic-steps:994.vd57e3ca_46d24
	h2-api:1.4.199
	javax-activation-api:1.2.0-5
	javax-mail-api:1.6.2-8
	jaxb:2.3.6-2
	jersey2-api:2.37-1
	jjwt-api:0.11.5-77.v646c772fddb_0
	jakarta-mail-api:2.0.1-2
	instance-identity:116.vf8f487400980
	jakarta-activation-api:2.0.1-2
	commons-text-api:1.10.0-27.vb_fa_3896786a_7
	aws-java-sdk-efs:1.12.287-357.vf82d85a_6eefd
	mina-sshd-api-common:2.9.1-44.v476733c11f82
	commons-lang3-api:3.12.0-36.vd97de6465d5b_`

	// plugins := `
	// jira:3.8`
	plugins = strings.ReplaceAll(plugins, " ", "")
	plugins = strings.ReplaceAll(plugins, "\t", "")
	// Split by ":"
	lines := strings.Split(plugins, "\n")
	for _, line := range lines {
		kv := strings.Split(line, ":")
		if len(kv) == 2 {
			r.AddJenkinsServerPlugin("jenkins-one", types.Project{
				Name:    kv[0],
				Owner:   "jenkinsci",
				Version: kv[1],
			})

		}
	}
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
