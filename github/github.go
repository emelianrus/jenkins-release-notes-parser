package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

type GitHub struct {
	Client http.Client

	RateLimit         int   // X-RateLimit-Limit | 60
	RateLimitRemaning int   // X-RateLimit-Remaining | 0
	RateLimitReset    int64 // X-RateLimit-Reset | 1679179139
	RateLimitUsed     int   // X-RateLimit-Used | 60
}

func Download(pluginName string) ([]types.GitHubReleaseNote, error) {
	github := GitHub{}
	fmt.Println("executed download goroutine " + pluginName)

	// Make a request to the API to get the release notes
	url := fmt.Sprintf("https://api.github.com/repos/jenkinsci/%s/releases", pluginName)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, errors.New("error making request")
	}

	// Loop until we get a successful response or hit the rate limit
	for {
		resp, err := github.Client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue // Try again
		}
		rateLimitRemaning, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
		github.RateLimitRemaning = rateLimitRemaning
		rateLimitUsed, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Used"))
		github.RateLimitUsed = rateLimitUsed
		rateLimit, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
		github.RateLimit = rateLimit

		if resp.StatusCode == http.StatusForbidden {
			// We hit the rate limit, so wait for the X-RateLimit-Reset header and try again
			fmt.Println("Hit rate limit, waiting...")

			resetTimestampStr := resp.Header.Get("X-RateLimit-Reset")
			resetTimestampInt, err := strconv.ParseInt(resetTimestampStr, 10, 64)
			github.RateLimitReset = resetTimestampInt

			resetTime := time.Unix(resetTimestampInt, 0).UTC()
			if err != nil {
				return nil, errors.New("converting date")
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
			fmt.Printf("API error: %s\n", resp.Status)
			return nil, fmt.Errorf("API error %s", resp.Status)
		}
		var releases []types.GitHubReleaseNote
		err = json.NewDecoder(resp.Body).Decode(&releases)
		if err != nil {
			// fmt.Println("error decoding github response")
			// log.Println(err)
			return nil, errors.New("error decoding github response")
		}

		fmt.Println("finished download goroutine " + pluginName)

		return releases, nil
	}
}
