package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetGitHubReleases(pluginName string) ([]GitHubReleaseNote, error) {
	fmt.Println("Downloading plugin from github " + pluginName)
	// Define cron job to run every hour
	// c := cron.New()
	// c.AddFunc("@hourly", func() {
	// log.Println("Fetching release notes from GitHub API...")

	// Make request to GitHub API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/jenkinsci/"+pluginName+"/releases", nil)
	if err != nil {
		// log.Printf("error in github request: %s\n", err)
		return nil, fmt.Errorf("error in github request: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		// log.Println(err)
		return nil, fmt.Errorf("error during making client.Do request: %s", err)
	}
	defer resp.Body.Close()

	// Check rate limit remaining
	// rateLimitRemaining, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// If rate limit is reached, wait for 1 hour
	// if rateLimitRemaining <= 0 {
	// 	log.Println("Rate limit reached. Waiting 1 hour...")
	// 	return
	// }

	// Decode response body into []GitHubReleaseNote
	var releases []GitHubReleaseNote
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		// fmt.Println("error decoding github response")
		// log.Println(err)
		return nil, fmt.Errorf("error decoding github response: %s", err)
	}
	return releases, nil
	// ownerName := "jenkinsci"
	// repoName := "plugin-installation-manager-tool"
	// // Cache releases in Redis for 1 hour

	// jsonData, err := json.Marshal(releases)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "lastUpdated"),
	// 	time.Now())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// log.Println("Releases cached in Redis")

	// tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))

	// Start cron job
	// c.Start()
	// c.Run()
	// Set up HTTP server
	// StartWeb(redisclient)
	// return releases
	// return nil
}
