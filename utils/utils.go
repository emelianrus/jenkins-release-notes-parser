package utils

import (
	"regexp"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// from slice of strings
func RemoveDuplicates(slice []string) []string {
	unique := make(map[string]bool)
	result := []string{}

	for _, val := range slice {
		if !unique[val] {
			unique[val] = true
			result = append(result, val)
		}
	}

	return result
}

func ConvertMarkdownToHtml(s string) string {
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

	resultHTML := string(markdown.Render(doc, renderer))
	return resultHTML
}

func ReplaceGitHubLinks(s string) string {
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
