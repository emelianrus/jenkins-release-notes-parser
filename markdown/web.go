package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"text/template"

	"github.com/go-redis/redis"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

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

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}
type Version struct {
	Version string
	Changes []string
}

type Product struct {
	Name     string
	Versions []Version
}

type PageData struct {
	Products []Product
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

func StartWeb(redisclient *redis.Client) {
	log.Println("Starting server")
	ownerName := "jenkinsci"
	repoName := "plugin-installation-manager-tool"

	// tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve releases from Redis
		key := fmt.Sprintf("github:%s:%s:2.10.0", ownerName, repoName)
		jsonData, err := redisclient.Get(key).Bytes()
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve releases from cache", http.StatusInternalServerError)
			return
		}

		var release Release
		err = json.Unmarshal(jsonData, &release)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to unmarshal releases from cache", http.StatusInternalServerError)
			return
		}

		// allKeys, _ := redisclient.Do("KEYS", "github:jenkinsci:plugin-installation-manager-tool:*").Result()

		// strings.Join(allKeys, ",")
		// strings.Split(a)

		// fmt.Println(res)

		html := convertMarkDownToHtml(release.Body)
		formatedHtml := replaceGitHubLinks(html)

		// htmlrendered += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", repoName, formatedHtml)
		// }
		// htmlrendered += "</tbody></table></body></html>"
		// fmt.Println(string(htmlrendered))
		// Write HTML response
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// w.Write([]byte(htmlrendered))

		products := []Product{
			{
				Name: repoName,
				Versions: []Version{
					{
						Version: "2.10.0",
						Changes: []string{formatedHtml},
					},
					{
						Version: "Version 2",
						Changes: []string{"Improved performance of feature A", "Added new option for feature B"},
					},
					{
						Version: "Version 3",
						Changes: []string{"Fixed bug in feature C", "Added new feature D", "Updated documentation"},
					},
				},
			},
			{
				Name: "Product 2",
				Versions: []Version{
					{
						Version: "Version 1",
						Changes: []string{"Added new feature X", "Fixed issue Y"},
					},
					{
						Version: "Version 2",
						Changes: []string{"Improved performance of feature X", "Added new option for feature Y", "Fixed bug in feature Z"},
					},
				},
			},
		}

		data := struct {
			Title      string
			Products   []Product
			ServerName string
		}{
			ServerName: "jenkins-one",
			Title:      "Plugin manager",
			Products:   products,
		}
		tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
