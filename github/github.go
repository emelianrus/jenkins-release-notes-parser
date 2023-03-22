package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
)

type GitHub struct {
	Initialized bool // first call to github will get all stats by quotes and rate limits
	Client      http.Client
	GitHubStats GitHubStats
}

type GitHubStats struct {
	RateLimit         int   // X-RateLimit-Limit | 60 | In units
	RateLimitRemaning int   // X-RateLimit-Remaining | 0 | In units
	RateLimitReset    int64 // X-RateLimit-Reset | 1679179139 | In seconds updated every ~hour
	RateLimitUsed     int   // X-RateLimit-Used | 60 | In units

	WaitSlotSeconds int // Seconds to reset RateLimit slots, if negative free to go
}

func NewGitHubClient() GitHub {
	return GitHub{}
}

func (g *GitHub) HasRequestSlot() bool {
	var reservedSlots = 5
	if g.GitHubStats.RateLimitRemaning == 0 && g.GitHubStats.RateLimit == 0 {
		fmt.Println("Can not get rate limit, github just created")
	}

	if g.GitHubStats.RateLimitRemaning < g.GitHubStats.RateLimit-reservedSlots {
		return true
	} else {
		return false
	}
}

// rate limit reset updated every ~hour
func (g *GitHub) SyncStats(resp *http.Response) {
	rateLimitRemaning, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	g.GitHubStats.RateLimitRemaning = rateLimitRemaning
	rateLimitUsed, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Used"))
	g.GitHubStats.RateLimitUsed = rateLimitUsed
	rateLimit, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
	g.GitHubStats.RateLimit = rateLimit
	resetTimestampStr := resp.Header.Get("X-RateLimit-Reset")
	resetTimestampInt, _ := strconv.ParseInt(resetTimestampStr, 10, 64)
	g.GitHubStats.RateLimitReset = resetTimestampInt

	nowEpoch := time.Now().Unix()
	timeToWait := resetTimestampInt - nowEpoch

	waitInt := int(timeToWait)

	g.GitHubStats.WaitSlotSeconds = waitInt

}

func (g *GitHub) updateWaitSlotSeconds() {
	nowEpoch := time.Now().Unix()
	timeToWait := g.GitHubStats.RateLimitReset - nowEpoch

	g.GitHubStats.WaitSlotSeconds = int(timeToWait)
	fmt.Printf("g.WaitSlotSeconds: %d\n", g.GitHubStats.WaitSlotSeconds)

}

func (g *GitHub) Download(pluginName string) ([]types.GitHubReleaseNote, error) {
	if g.Initialized {
		g.updateWaitSlotSeconds()
	}

	// we need to add suffix to plugins it differs plugin name and github project
	if !strings.HasSuffix(pluginName, "-plugin") {
		pluginName = pluginName + "-plugin"
	}

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
		resp, err := g.Client.Do(req)
		g.SyncStats(resp)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue // Try again
		}

		if resp.StatusCode == http.StatusForbidden {
			// We hit the rate limit, so wait for the X-RateLimit-Reset header and try again
			fmt.Printf("Rate limit reached, waiting until %d seconds\n", g.GitHubStats.WaitSlotSeconds)
			time.Sleep(time.Second * time.Duration(g.GitHubStats.WaitSlotSeconds))

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
