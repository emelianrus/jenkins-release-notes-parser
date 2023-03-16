package main

func main() {
	redisclient := NewRedisClient()

	// TODO: this data is temporary, used during development
	// should be replaced by client call
	redisclient.addJenkinsServer("jenkins-two", "2.3233.1")
	redisclient.addJenkinsServer("jenkins-one", "2.3233.2")

	redisclient.addJenkinsServerPlugin("jenkins-one", JenkinsPlugin{
		Name:    "plugin-installation-manager-tool",
		Version: "2.10.0",
	})
	redisclient.addJenkinsServerPlugin("jenkins-two", JenkinsPlugin{
		Name:    "plugin-installation-manager-tool",
		Version: "2.10.0",
	})
	redisclient.addJenkinsServerPlugin("jenkins-one", JenkinsPlugin{
		Name:    "okhttp-api-plugin",
		Version: "4.9.3-108.v0feda04578cf",
	})

	StartWeb(redisclient)
}
