package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

func Download(pluginName string) ([]types.GitHubReleaseNote, error) {
	fmt.Println("executed download goroutine " + pluginName)
	client := http.Client{}

	// Make a request to the API to get the release notes
	url := fmt.Sprintf("https://api.github.com/repos/jenkinsci/%s/releases", pluginName)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, nil
	}

	// Loop until we get a successful response or hit the rate limit
	for {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue // Try again
		}

		if resp.StatusCode == http.StatusForbidden {
			// We hit the rate limit, so wait for the X-RateLimit-Reset header and try again
			fmt.Println("Hit rate limit, waiting...")

			resetTimestampStr := resp.Header.Get("X-RateLimit-Reset")
			resetTimestampInt, err := strconv.ParseInt(resetTimestampStr, 10, 64)

			resetTime := time.Unix(resetTimestampInt, 0).UTC()
			if err != nil {
				panic(err)
			}
			nowEpoch := time.Now().Unix()
			timeToWait := resetTimestampInt - nowEpoch

			waitInt := int(timeToWait)
			fmt.Printf("Rate limit reached, waiting until %d seconds, until %s...\n", waitInt, resetTime)
			time.Sleep(time.Second * time.Duration(waitInt))

			continue // Try again
		}

		if resp.StatusCode != http.StatusOK {
			// The API returned an error, so print the status code and message and exit
			fmt.Printf("API error: %s - %s", resp.Status, resp.Body)
			return nil, nil
		}
		var releases []types.GitHubReleaseNote
		err = json.NewDecoder(resp.Body).Decode(&releases)
		if err != nil {
			// fmt.Println("error decoding github response")
			// log.Println(err)
			return nil, nil
		}

		// The request was successful, so print the release notes and exit
		// fmt.Println("Release notes:")
		// fmt.Println(releases)

		fmt.Println("finished download goroutine " + pluginName)

		return releases, nil
	}
}
