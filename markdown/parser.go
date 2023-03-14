package main

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

// From redis
type Versions []string

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

var ownerName = "jenkinsci"

func GetGitHubReleases(pluginName string, redisclient *Redis) {

	// Define cron job to run every hour
	// c := cron.New()
	// c.AddFunc("@hourly", func() {
	// log.Println("Fetching release notes from GitHub API...")

	// Make request to GitHub API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/jenkinsci/"+pluginName+"/releases", nil)
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
	// ownerName := "jenkinsci"
	// repoName := "plugin-installation-manager-tool"
	// // Cache releases in Redis for 1 hour

	// jsonData, err := json.Marshal(releases)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", ownerName, pluginName, "lastUpdated"),
		time.Now().Unix())
	if err != nil {
		log.Println(err)
		return
	}

	versions := Versions{}

	for _, release := range releases {
		versions = append(versions, release.Name)
		key := fmt.Sprintf("github:%s:%s:%s", ownerName, pluginName, release.Name)
		// 0 time.Hour
		jsonData, err := json.Marshal(release)
		if err != nil {
			log.Println(err)
			return
		}
		err = redisclient.Set(key, jsonData)
		if err != nil {
			log.Println(err)
			return
		}
	}
	jsonVersions, _ := json.Marshal(versions)
	err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", ownerName, pluginName, "versions"),
		jsonVersions)
	if err != nil {
		log.Println(err)
		return
	}

	// log.Println("Releases cached in Redis")

	// tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))

	// Start cron job
	// c.Start()
	// c.Run()
	// Set up HTTP server
	// StartWeb(redisclient)
	// return releases
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

// get jenkins server (hardcoded)
// get plugins from jenikins server
// check cache for plugin by versions file
// construct pageData
func getPlugins(redisclient *Redis) PluginPage {

	// ownerName := "jenkinsci"
	// repoName := "plugin-installation-manager-tool"

	serverJson, _ := redisclient.GetJenkinsServers()

	var jenkinsServer JenkinsServer
	err := json.Unmarshal(serverJson, &jenkinsServer)
	if err != nil {
		fmt.Println("can not unmarshal jenkins server")
	}

	fmt.Println(jenkinsServer.Name)
	fmt.Println(jenkinsServer.Plugins)
	products := []Product{}
	for _, plugin := range jenkinsServer.Plugins {

		pluginVersionsJson, err := redisclient.GetPlugin(fmt.Sprintf("github:jenkinsci:%s:versions", plugin.Name))
		if err != nil {
			// plugin doesnt exist
			fmt.Println("versions file is not exist")
			fmt.Println(err)
			GetGitHubReleases(plugin.Name, redisclient)
			// plugin = redisclient.Get()
		}

		// Assume we hit redis cache
		var versions Versions
		err = json.Unmarshal(pluginVersionsJson, &versions)
		if err != nil {
			log.Println(err)
			// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
			return PluginPage{}
		}

		var convertedVersions []Version
		for _, version := range versions {

			pluginJson, _ := redisclient.Get(fmt.Sprintf("github:jenkinsci:%s:%s", plugin.Name, version)).Bytes()

			var releaseNote Release
			err = json.Unmarshal(pluginJson, &releaseNote)

			if err != nil {
				log.Println(err)
				// http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
				return PluginPage{}
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

	return PluginPage{
		Title:      "Plugin manager",
		ServerName: "jenkins-one",
		Products:   products,
	}
}
