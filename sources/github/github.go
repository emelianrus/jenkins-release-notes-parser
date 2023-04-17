package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
	"github.com/sirupsen/logrus"
)

type GitHub struct {
	Initialized bool // first call to github will get all stats by quotes and rate limits
	Client      http.Client
	GitHubStats GitHubStats
}

type GitHubStats struct {
	RateLimit         int   // X-RateLimit-Limit 		| 60 			| In units
	RateLimitRemaning int   // X-RateLimit-Remaining 	| 0 			| In units
	RateLimitReset    int64 // X-RateLimit-Reset 		| 1679179139 	| In seconds updated every ~hour
	RateLimitUsed     int   // X-RateLimit-Used 		| 60			| In units

	WaitSlotSeconds int // Seconds to reset RateLimit slots, if negative free to go
}

// response from github getted from release path
type gitHubReleaseNote struct {
	Name      string `json:"name"`       // Version
	TagName   string `json:"tag_name"`   // tag name
	Body      string `json:"body"`       // this is markdown formated text of release note
	CreatedAt string `json:"created_at"` //
}

func NewGitHubClient() GitHub {
	return GitHub{}
}

func (g *GitHub) HasRequestSlot() bool {
	var reservedSlots = 5
	if g.GitHubStats.RateLimitRemaning == 0 && g.GitHubStats.RateLimit == 0 {
		logrus.Infoln("Can not get rate limit, github just created")
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
	logrus.Infof("g.WaitSlotSeconds: %d\n", g.GitHubStats.WaitSlotSeconds)

}

func (g *GitHub) waitUntilNextSlotAvailable() {
	logrus.Infof("Rate limit reached, waiting until %d seconds\n", g.GitHubStats.WaitSlotSeconds)
	time.Sleep(time.Second * time.Duration(g.GitHubStats.WaitSlotSeconds))
}

func (g *GitHub) Download(projectName string) ([]types.ReleaseNote, error) {
	releaseNotes := []types.ReleaseNote{}

	if g.Initialized {
		g.updateWaitSlotSeconds()
	}

	// use only half of the slots
	if g.GitHubStats.RateLimitUsed > 30 {
		g.waitUntilNextSlotAvailable()
	}

	// TODO: how to handle this with no rely on jenkins specific projects
	// we need to add suffix to plugins it differs plugin name and github project
	if !strings.HasSuffix(projectName, "-plugin") {
		logrus.Warnf("project %s doesn't have 'Name'\n", projectName)
		projectName = projectName + "-plugin"
	}

	logrus.Infoln("executed download goroutine " + projectName)

	// Make a request to the API to get the release notes
	url := fmt.Sprintf("https://api.github.com/repos/jenkinsci/%s/releases", projectName)
	logrus.Infoln(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorln("Error creating request:", err)
		return nil, errors.New("error making request")
	}

	// Loop until we get a successful response or hit the rate limit
	for {
		resp, err := g.Client.Do(req)
		g.SyncStats(resp)
		logrus.Infof("RateLimitUsed: %d\n", g.GitHubStats.RateLimitUsed)
		if err != nil {
			logrus.Errorln("Error making request:", err)
			continue // Try again
		}

		if resp.StatusCode == http.StatusForbidden {
			// We hit the rate limit, so wait for the X-RateLimit-Reset header and try again
			// fmt.Printf("Rate limit reached, waiting until %d seconds\n", g.GitHubStats.WaitSlotSeconds)
			// time.Sleep(time.Second * time.Duration(g.GitHubStats.WaitSlotSeconds))
			g.waitUntilNextSlotAvailable()
			continue // Try again
		}

		if resp.StatusCode != http.StatusOK {
			// The API returned an error, so print the status code and message and exit
			logrus.Errorf("API error: %s\n", resp.Status)
			return nil, fmt.Errorf("API error %s", resp.Status)
		}
		var releases []gitHubReleaseNote
		err = json.NewDecoder(resp.Body).Decode(&releases)
		if err != nil {
			// fmt.Println("error decoding github response")
			// log.Println(err)
			return nil, errors.New("error decoding github response")
		}

		logrus.Infoln("finished download goroutine " + projectName)

		for _, release := range releases {
			releaseNotes = append(releaseNotes, types.ReleaseNote{
				Name: release.Name,
				Tag:  release.TagName,
				BodyHTML: string(template.HTML(
					utils.ReplaceGitHubLinks(
						utils.ConvertMarkdownToHtml(release.Body)))),
				CreatedAt: release.CreatedAt,
			})
		}
		return releaseNotes, nil
	}
}
