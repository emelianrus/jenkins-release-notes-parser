package main

import (
	"fmt"
	"regexp"
)

// (#526)
// https://github.com/jenkinsci/plugin-installation-manager-tool/pull/526

// @timja
// https://github.com/timja

var str = `<li>Bump checkstyle from 8.37 to 8.38 (#235) @dependabot</li>

// <li>Bump checkstyle from 8.37 to 8.38 (#236) @dependabot</li>
// <li>Bump checkstyle from 8.37 to 8.38 (#2361) @dependabot</li>
// <li>Bump checkstyle from 8.37 to 8.38 (#2) @dependabot</li>`

func main() {

	// Replace the number with a link
	re := regexp.MustCompile(`#(\d+)`)
	output1 := re.ReplaceAllString(str, `<a href="https://github.com/jenkinsci/plugin-installation-manager-tool/pull/$1">#$1</a>`)

	// Replace dynamic strings with links
	re2 := regexp.MustCompile(`@([-\w\d]+)`)
	output2 := re2.ReplaceAllString(output1, `<a href="https://github.com/$1">@$1</a>`)

	fmt.Println(output2)
}
