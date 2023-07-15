package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/emelianrus/jenkins-release-notes-parser/utils"
	"github.com/sirupsen/logrus"
)

type GitHub struct {
	Initialized   bool // first call to github will get all stats by quotes and rate limits
	Client        http.Client
	GitHubStats   GitHubStats
	PersonalToken string
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
	HtmlUrl   string `json:"html_url"`
}

func NewGitHubClient() GitHub {
	gh := GitHub{}

	personalToken := os.Getenv("GITHUB_PERSONAL_TOKEN")

	if personalToken != "" {
		gh.SetPersonalToken(personalToken)
		logrus.Infoln("Using personal token for github connection")
	}
	return gh
}

func (g *GitHub) SetPersonalToken(token string) {
	g.PersonalToken = token
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

	g.GitHubStats.WaitSlotSeconds = int(timeToWait)

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
	logrus.Infof("[GITHUB][Download] project %s\n", projectName)
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
		logrus.Warnf("[GITHUB][Download] project %s doesn't have 'Name'\n", projectName)
		projectName = projectName + "-plugin"
	}

	logrus.Infoln("[GITHUB][Download] executed download goroutine " + projectName)

	// Make a request to the API to get the release notes
	// TODO: github can return only 1000 release notes
	// need to add 'page=1' to 'page=10'
	// currently can not return more then 100 releases from latest version
	url := fmt.Sprintf("https://api.github.com/repos/jenkinsci/%s/releases?per_page=100", projectName)
	logrus.Infoln(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorln("[GITHUB][Download] Error creating request:", err)
		return nil, errors.New("error making request")
	}

	if g.PersonalToken != "" {
		// Set the request headers
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Set("Authorization", "Bearer "+g.PersonalToken)
	}

	// Loop until we get a successful response or hit the rate limit
	for {
		resp, err := g.Client.Do(req)
		g.SyncStats(resp)
		logrus.Infof("%+v\n", g.GitHubStats)
		logrus.Infof("RateLimitUsed: %d from %d\n", g.GitHubStats.RateLimitUsed, g.GitHubStats.RateLimit)
		if err != nil {
			logrus.Errorln("Error making request:", err)
			continue // Try again
		}

		if resp.StatusCode == http.StatusUnauthorized {
			logrus.Errorln("GITHUB_PERSONAL_TOKEN is not corrent or expired, will use public api")
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
				Name:    release.Name,
				Tag:     release.TagName,
				HTMLURL: release.HtmlUrl,
				BodyHTML: string(template.HTML(
					utils.ReplaceGitHubLinks(
						utils.ConvertMarkdownToHtml(release.Body)))),
				CreatedAt: release.CreatedAt,
			})
		}
		return releaseNotes, nil
	}
}
