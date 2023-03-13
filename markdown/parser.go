package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type Release struct {
	Name      string `json:"name"` // Version
	Body      string `json:"body"` // this is markdown formated text of release note
	CreatedAt string `json:"created_at"`
}

// represent plugin
type Plugin struct {
	Name         string
	ReleaseNotes []Release
}

// represent jenkins server
type Server struct {
	Plugins []Plugin
}
type Versions []string

// func getFromGitHub() {

// }

// func geatherPlugins() {

// }

func Parser(redisclient *redis.Client) {

	// Define cron job to run every hour
	// c := cron.New()
	// c.AddFunc("@hourly", func() {
	// log.Println("Fetching release notes from GitHub API...")

	// Make request to GitHub API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/jenkinsci/plugin-installation-manager-tool/releases", nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
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

	// Decode response body into []Release
	var releases []Release
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		log.Println(err)
		return
	}
	ownerName := "jenkinsci"
	repoName := "plugin-installation-manager-tool"
	// Cache releases in Redis for 1 hour

	// jsonData, err := json.Marshal(releases)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", ownerName, repoName, "lastUpdated"),
		time.Now().Unix(), 0).Err()
	if err != nil {
		log.Println(err)
		return
	}

	versions := Versions{}

	for _, release := range releases {
		versions = append(versions, release.Name)
		key := fmt.Sprintf("github:%s:%s:%s", ownerName, repoName, release.Name)
		// 0 time.Hour
		jsonData, err := json.Marshal(release)
		if err != nil {
			log.Println(err)
			return
		}
		err = redisclient.Set(key, jsonData, 0).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}
	jsonVersions, _ := json.Marshal(versions)
	err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", ownerName, repoName, "versions"),
		jsonVersions, 0).Err()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Releases cached in Redis")

	// tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))

	// Start cron job
	// c.Start()
	// c.Run()
	// Set up HTTP server
	StartWeb(redisclient)

}
