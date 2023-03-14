package main

/*

we have from github
  all versions release notes

  1) check if we have plugin data in cache?
  if (exist) {
	use
  } else {
	get release from github
	add to cache
  }

*/
import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type GitHubReleaseNote struct {
	Name      string `json:"name"` // Version
	Body      string `json:"body"` // this is markdown formated text of release note
	CreatedAt string `json:"created_at"`
}

// represent plugin
type Plugin struct {
	Name         string
	ReleaseNotes []GitHubReleaseNote
}

// represent jenkins server
type Server struct {
	Plugins []Plugin
}

// From redis
type Versions []string

// HTML start
// part of html responce
type Version struct {
	Version string
	Changes template.HTML
}

type Product struct {
	Name             string
	Versions         []Version
	InstalledVersion string
	LastUpdated      string // TODO
}
type PluginPage struct {
	Title      string
	Products   []Product
	ServerName string
}

// HTML end
var ownerName = "jenkinsci"

func GetGitHubReleases(pluginName string, redisclient *Redis) ([]GitHubReleaseNote, error) {
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

func convertMarkDownToHtml(s string) string {
	md := []byte(s)
	// always normalize newlines, this library only supports Unix LF newlines
	md = markdown.NormalizeNewlines(md)
	// create markdown parser
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// parse markdown into AST tree
	doc := p.Parse(md)
	// create HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := string(markdown.Render(doc, renderer))
	return html
}

func replaceGitHubLinks(s string) string {
	// (#526)
	// https://github.com/jenkinsci/plugin-installation-manager-tool/pull/526

	// @timja
	// https://github.com/timja

	re := regexp.MustCompile(`#(\d+)`)
	output1 := re.ReplaceAllString(s, `<a href="https://github.com/jenkinsci/plugin-installation-manager-tool/pull/$1">#$1</a>`)

	// Replace dynamic strings with links
	re2 := regexp.MustCompile(`@([-\w\d]+)`)
	output2 := re2.ReplaceAllString(output1, `<a href="https://github.com/$1">@$1</a>`)
	return output2
}

func saveReleaseNotesToDB(redisclient *Redis, releases []GitHubReleaseNote, pluginName string) error {

	currentTime := time.Now()
	formattedTime := currentTime.Format("02 January 2006 15:04")
	err := redisclient.SetLastUpdatedTime(pluginName, formattedTime)
	if err != nil {
		// fmt.Println(err)
		// fmt.Println("Can not set updated time")
		return fmt.Errorf("error setting lastUpdate time in get github release: %s", err)
	}
	versions := Versions{}

	for _, release := range releases {
		versions = append(versions, release.Name)
		key := fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, release.Name)
		// 0 time.Hour
		jsonData, err := json.Marshal(release)
		if err != nil {
			// log.Println(err)
			return fmt.Errorf("error Marshal release: %s", err)
		}
		err = redisclient.Set(key, jsonData)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error setting release: %s", err)
		}
	}

	jsonVersions, _ := json.Marshal(versions)
	err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", "jenkinsci", pluginName, "versions"),
		jsonVersions)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error setting version for release: %s", err)
	}
	return nil
}

// get jenkins server (hardcoded)
// get plugins from jenikins server
// check cache for plugin by versions file
// construct pageData
func getPluginsForPageData(redisclient *Redis, jenkinsServer JenkinsServer) PluginPage {
	// default page data
	plPage := PluginPage{
		Title:      "Plugin manager",
		ServerName: "jenkins-one",
		Products:   nil,
	}

	products := []Product{}

	for _, plugin := range jenkinsServer.Plugins {

		pluginVersionsJson, err := redisclient.GetPluginVersions(plugin.Name)
		if err != nil {
			fmt.Println("versions file doesn't exist in redis cache for " + plugin.Name)
			fmt.Println(err)
			releases, err := GetGitHubReleases(plugin.Name, redisclient)
			if err != nil {
				fmt.Println("Failed to get releases from github")
			}
			err = saveReleaseNotesToDB(redisclient, releases, plugin.Name)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Failed to save release notes to db")
			}

			pluginVersionsJson, err = redisclient.GetPluginVersions(plugin.Name)
			if err != nil {
				fmt.Println(err)
				fmt.Println("2nd attempt to GetPluginVersions failed")
				// return web page with default values
				return plPage
			}
		}

		// Assume we hit redis cache
		var versions Versions
		err = json.Unmarshal(pluginVersionsJson, &versions)
		if err != nil {
			log.Println(err)
			// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
			return plPage
		}

		var convertedVersions []Version
		for _, version := range versions {

			releaseNote, err := redisclient.GetPluginWithVersion(plugin.Name, version)
			if err != nil {
				log.Println(err)
				// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
				return plPage
			}

			convertedVersions = append(convertedVersions, Version{
				Version: version,
				Changes: template.HTML(replaceGitHubLinks(convertMarkDownToHtml(releaseNote.Body))),
			})
		}

		lastUpdated, _ := redisclient.Get(fmt.Sprintf("github:%s:%s:%s", ownerName, plugin.Name, "lastUpdated")).Bytes()

		products = append(products,
			Product{
				Name:             plugin.Name,
				Versions:         convertedVersions,
				InstalledVersion: plugin.Version,
				LastUpdated:      string(lastUpdated),
			},
		)
	}
	plPage.Products = products
	return plPage
}
