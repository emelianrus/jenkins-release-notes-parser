package main

func main() {
	redisclient := NewRedisClient()

	// init jenkins server
	// jenkinsServer1 := JenkinsServer{
	// 	Name: "jenkins-one",
	// 	Core: "2.3233.2",
	// }

	// init jenkins server
	// jenkinsServer2 := JenkinsServer{
	// 	Name: "jenkins-two",
	// 	Core: "2.3233.2",
	// }
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
	// redisclient.removeJenkinsServerPlugin("jenkins-one", "newplugin")
	// redisclient.changeJenkinServerPluginVersion("jenkins-one", "newplugin", "111.111.11")
	// res := redisclient.getJenkinsPlugins("jenkins-one")
	// fmt.Println(res)

	// jsonData, err := json.Marshal(js)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// // write jenkins server json
	// err = redisclient.Set("servers:jenkins-one:plugins", jsonData)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// getPlugins(redisclient)
	// // GetPluginFromGitHub(redisclient)

	// redisclient.getJenkinsServers()
	StartWeb(redisclient)

	// plugin, err := redisclient.GetPlugin("github:jenkinsci:plugin-installation-manager-tool:versions")

	// fmt.Println(err)
	// fmt.Println(plugin)
}
