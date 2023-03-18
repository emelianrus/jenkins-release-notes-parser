package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

var (
	serviceMutex sync.Mutex
)

// used as go StartQueue()
// can be executed by button from UI so we need to be sure running only one instance at once
func StartQueue(redisclient *db.Redis) {
	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	for {
		for _, pluginName := range redisclient.GetAllPluginsFromServers() {
			// TODO: error api 404
			ghReleaseNotes := download(pluginName)
			redisclient.SaveReleaseNotesToDB(ghReleaseNotes, pluginName)
		}
		fmt.Println("sleep an hour")
		time.Sleep(time.Hour)
	}
}

func download(pluginName string) []types.GitHubReleaseNote {
	fmt.Println("executed download goroutine " + pluginName)
	// Set up an HTTP client with a rate limiter
	rate := time.Second / 15 // Allow 30 requests per second
	client := http.Client{}
	ticker := time.NewTicker(rate)
	defer ticker.Stop()

	// Make a request to the API to get the release notes
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/jenkinsci/%s/releases", pluginName), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Loop until we get a successful response or hit the rate limit
	for {
		<-ticker.C // Wait for the next available request slot

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue // Try again
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			// We hit the rate limit, so wait for the Retry-After header and try again
			fmt.Println("Hit rate limit, waiting...")
			retryAfter, err := time.ParseDuration(resp.Header.Get("Retry-After") + "s")
			if err != nil {
				fmt.Println("Error parsing Retry-After header:", err)
				return nil
			}
			fmt.Println("next try at")
			fmt.Println(retryAfter)
			// TODO: set retryAfter to redis "github:retryAfter"
			ticker.Reset(retryAfter) // Wait for the specified duration before trying again
			continue                 // Try again
		}

		if resp.StatusCode != http.StatusOK {
			// The API returned an error, so print the status code and message and exit
			fmt.Printf("API error: %s - %s", resp.Status, resp.Body)
			return nil
		}
		var releases []types.GitHubReleaseNote
		err = json.NewDecoder(resp.Body).Decode(&releases)
		if err != nil {
			// fmt.Println("error decoding github response")
			// log.Println(err)
			return nil
		}

		// The request was successful, so print the release notes and exit
		// fmt.Println("Release notes:")
		// fmt.Println(releases)

		fmt.Println("finished download goroutine " + pluginName)

		return releases
	}
}
