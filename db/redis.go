package db

import (
	"fmt"
	"strings"

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
	// r.AddJenkinsServer("jenkins-two", "2.3233.1")
	r.AddJenkinsServer("jenkins-one", "2.3233.2")

	plugins := `
	ace-editor:1.1
	ansicolor:1.0.2
	ant:481.v7b_09e538fcca
	antisamy-markup-formatter:2.7
	apache-httpcomponents-client-4-api:4.5.13-138.v4e7d9a_7b_a_e61
	artifactory:3.17.1
	audit-trail:3.11
	authentication-tokens:1.4
	aws-credentials:191.vcb_f183ce58b_9
	aws-java-sdk:1.12.287-357.vf82d85a_6eefd
	aws-lambda:0.5.10
	basic-branch-build-strategies:1.3.2
	blueocean-autofavorite:1.2.5
	blueocean-bitbucket-pipeline:1.25.8
	blueocean-commons:1.25.8
	blueocean-config:1.25.8
	blueocean-core-js:1.25.8
	blueocean-dashboard:1.25.8
	blueocean-display-url:2.4.1
	blueocean-events:1.25.8
	blueocean-git-pipeline:1.25.8
	blueocean-github-pipeline:1.25.8
	blueocean-i18n:1.25.8
	blueocean-jira:1.25.8
	blueocean-jwt:1.25.8
	blueocean-personalization:1.25.8
	blueocean-pipeline-api-impl:1.25.8
	blueocean-pipeline-editor:1.25.8
	blueocean-pipeline-scm-api:1.25.8
	blueocean-rest-impl:1.25.8
	blueocean-rest:1.25.8
	blueocean-web:1.25.8
	blueocean:1.25.8
	bouncycastle-api:2.26
	branch-api:2.1046.v0ca_37783ecc5
	build-name-setter:2.2.0
	build-timeout:1.24
	cloudbees-bitbucket-branch-source:791.vb_eea_a_476405b
	cloudbees-disk-usage-simple:178.v1a_4d2f6359a_8
	cloudbees-folder:6.758.vfd75d09eea_a_1
	command-launcher:90.v669d7ccb_7c31
	compress-artifacts:1.10
	conditional-buildstep:1.4.2
	config-file-provider:3.11.1
	configuration-as-code:1559.v38a_b_2e3b_6b_b_7
	copyartifact:1.46.3
	credentials-binding:523.vd859a_4b_122e6
	credentials:1143.vb_e8b_b_ceee347
	display-url-api:2.3.6
	docker-commons:1.21
	docker-workflow:521.v1a_a_dd2073b_2e
	durable-task:501.ve5d4fc08b0be
	dynamic-axis:1.0.3
	email-ext:2.92
	embeddable-build-status:255.va_d2370ee8fde
	favorite:2.4.1
	file-parameters:205.vf6ce13b_e5dee
	git-client:3.12.1
	git-server:99.va_0826a_b_cdfa_d
	git:4.12.1
	github-api:1.303-400.v35c2d8258028
	github-branch-source:1695.v88de84e9f6b_9
	github:1.35.0
	gitlab-api:5.0.1-78.v47a_45b_9f78b_7
	gitlab-branch-source:640.v7101b_1c0def9
	gitlab-logo:1.0.5
	gitlab-oauth:1.16
	gitlab-plugin:1.5.35
	gradle:1.40
	handlebars:3.0.8
	handy-uri-templates-2-api:2.1.8-22.v77d5b_75e6953
	hashicorp-vault-plugin:356.ved18810a_b_828
	htmlpublisher:1.31
	ivy:2.2
	jackson2-api:2.13.4-293.vee957901b_6fb
	javadoc:226.v71211feb_e7e9
	jdk-tool:55.v1b_32b_6ca_f9ca
	jenkins-design-language:1.25.8
	jira:3.8
	job-dsl:1.81
	job-import-plugin:3.5
	jobConfigHistory:1176.v1b_4290db_41a_5
	jquery-detached:1.2.1
	jquery:1.12.4-1
	jsch:0.1.55.61.va_e9ee26616e7
	junit:1153.v1c24f1a_d2553
	kubernetes-cli:1.10.3
	kubernetes:3724.v0920c1e0ec69
	lockable-resources:2.18
	mailer:438.v02c7f0a_12fa_4
	matrix-auth:3.1.5
	matrix-project:785.v06b_7f47b_c631
	maven-plugin:3.20
	mercurial:1251.va_b_121f184902
	metrics:4.2.10-389.v93143621b_050
	momentjs:1.1.1
	monitoring:1.91.0
	nodejs:1.5.1
	nodelabelparameter:1.11.0
	opentelemetry:2.9.2
	parameterized-scheduler:1.1
	parameterized-trigger:2.45
	Parameterized-Remote-Trigger:3.1.6.3
	pipeline-aws:1.43
	pipeline-build-step:2.18
	pipeline-graph-analysis:195.v5812d95a_a_2f9
	pipeline-input-step:451.vf1a_a_4f405289
	pipeline-maven:1205.vceea_7b_972817
	pipeline-milestone-step:101.vd572fef9d926
	pipeline-model-api:2.2114.v2654ca_721309
	pipeline-model-definition:2.2114.v2654ca_721309
	pipeline-model-extensions:2.2114.v2654ca_721309
	pipeline-rest-api:2.26
	pipeline-stage-step:296.v5f6908f017a_5
	pipeline-stage-tags-metadata:2.2114.v2654ca_721309
	pipeline-stage-view:2.26
	pipeline-utility-steps:2.13.0
	plain-credentials:139.ved2b_9cf7587b
	pubsub-light:1.17
	resource-disposer:0.20
	run-condition:1.5
	saml:4.354.vdc8c005cda_34
	sauce-ondemand:1.207
	scm-api:621.vda_a_b_055e58f7
	script-security:1183.v774b_0b_0a_a_451
	selenium-axis:0.0.6
	slack:625.va_eeb_b_168ffb_0
	sse-gateway:1.26
	ssh-agent:295.v9ca_a_1c7cc3a_a_
	ssh-credentials:305.v8f4381501156
	ssh-slaves:2.846.v1b_70190624f5
	structs:324.va_f5d6774f3a_d
	swarm:3.34
	throttle-concurrents:2.9
	timestamper:1.20
	token-macro:308.v4f2b_ed62b_b_16
	trilead-api:2.72.v2a_3236754f73
	uno-choice:2.6.4
	variant:59.vf075fe829ccb
	warnings-ng:9.20.1
	windows-slaves:1.8.1
	workflow-aggregator:590.v6a_d052e5a_a_b_5
	workflow-api:1200.v8005c684b_a_c6
	workflow-basic-steps:994.vd57e3ca_46d24
	workflow-cps-global-lib:588.v576c103a_ff86
	workflow-cps:2802.v5ea_628154b_c2
	workflow-durable-task-step:1206.v8a_d5f86e336b
	workflow-job:1239.v71b_b_a_124a_725
	workflow-multibranch:716.vc692a_e52371b_
	workflow-scm-step:400.v6b_89a_1317c9a_
	workflow-step-api:639.v6eca_cd8c04a_a_
	workflow-support:838.va_3a_087b_4055b
	ws-cleanup:0.43
	analysis-model-api:10.17.0
	aws-java-sdk-ec2:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-minimal:1.12.287-357.vf82d85a_6eefd
	bootstrap4-api:4.6.0-5
	bootstrap5-api:5.2.1-3
	caffeine-api:2.9.3-65.v6a_47d0f4d1fe
	checks-api:1.7.5
	data-tables-api:1.12.1-4
	echarts-api:5.4.0-1
	extended-read-permission:3.2
	font-awesome-api:6.2.0-3
	forensics-api:1.16.0
	google-container-registry-auth:0.3
	google-oauth-plugin:1.0.7
	google-play-android-publisher:4.2
	h2-api:1.4.199
	javax-activation-api:1.2.0-5
	javax-mail-api:1.6.2-8
	jaxb:2.3.6-2
	jersey2-api:2.37-1
	jjwt-api:0.11.5-77.v646c772fddb_0
	jquery3-api:3.6.1-2
	kubernetes-client-api:5.12.2-193.v26a_6078f65a_9
	kubernetes-credentials:0.9.0
	node-iterator-api:49.v58a_8b_35f8363
	oauth-credentials:0.5
	okhttp-api:4.9.3-108.v0feda04578cf
	pipeline-groovy-lib:612.v84da_9c54906d
	plugin-util-api:2.18.0
	popper-api:1.16.1-3
	popper2-api:2.11.6-2
	prism-api:1.29.0-1
	rebuild:1.34
	snakeyaml-api:1.32-86.ve3f030a_75631
	sonar:2.14
	sshd:3.249.v2dc2ea_416e33
	aws-java-sdk-sns:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-elasticbeanstalk:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-ssm:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-cloudformation:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-ecs:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-iam:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-sqs:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-ecr:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-logs:1.12.287-357.vf82d85a_6eefd
	aws-java-sdk-codebuild:1.12.287-357.vf82d85a_6eefd
	ionicons-api:31.v4757b_6987003
	mina-sshd-api-core:2.9.1-44.v476733c11f82
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
			r.AddJenkinsServerPlugin("jenkins-one", types.JenkinsPlugin{
				Name:    kv[0],
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
