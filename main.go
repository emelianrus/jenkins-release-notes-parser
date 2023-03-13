// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	git "github.com/emelianrus/jenkins-release-notes-parser/pkg/releasenotes/git"
// 	releaseParser "github.com/emelianrus/jenkins-release-notes-parser/pkg/releasenotes/releaseParser"
// )

// func main() {
// 	byt, _ := ioutil.ReadFile("diff.out")
// 	diff, _ := git.Parse(string(byt))

// 	project, _ := releaseParser.GetChangedVersions(diff)
// 	for _, f := range project.Files {

// 		var releaseNotes string
// 		for _, p := range f.Plugins {
// 			releases, err := releaseParser.GetReleases(p.Name)
// 			if err != nil {
// 				fmt.Println(err)
// 				break
// 			}
// 			_, result := releaseParser.GetReleaseNotes(releases, p)

// 			releaseNotes += result
// 		}

// 		err := os.WriteFile(f.Name+"_RELEASE_NOTES.md", []byte(releaseNotes), 0644)

// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 	}
// }

package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"fmt"
	"os"
)

var mdStr = `
# aws-java-sdk:aws-java-sdk-1.11.955
<!-- Optional: add a release summary here -->
## 📦 Dependency updates

* Bump plugin from 4.15 to 4.16 (#365) @dependabot
* Bump aws-java-sdk from 1.11.931 to 1.11.955 (#375) @dependabot

 ---
# docker-commons:docker-commons-1.19
<!-- Optional: add a release summary here -->
## 🐛 Bug fixes

* [JENKINS-67572](https://issues.jenkins.io/browse/JENKINS-67572) - Allow docker digest in image names (#93) @viceice

 ---
# docker-commons:docker-commons-1.18
<!-- Optional: add a release summary here -->
* [SECURITY-1878](https://www.jenkins.io/security/advisory/2022-01-12/)

 ---
# credentials-binding:credentials-binding-1.27
<!-- Optional: add a release summary here -->
## 🐛 Bug fixes

* Do not suggest 'passphraseVariable: '', usernameVariable: ''' in snippet generator (#144) @jglick
* [JENKINS-64361](https://issues.jenkins.io/browse/JENKINS-64361) - Make fix for [JENKINS-44860](https://issues.jenkins.io/browse/JENKINS-44860) - apply to Pipeline step arguments as well (take II) (#143) @jglick

## 📦 Dependency updates

* Bump git-changelist-maven-extension from 1.0-beta-7 to 1.2 (#142) @dependabot
* Bump bom-2.235.x from 872.v03c18fa35487 to 887.vae9c8ac09ff7 (#140) @dependabot

 ---
`

func main() {
	md := []byte(mdStr)
	// always normalize newlines, this library only supports Unix LF newlines
	md = markdown.NormalizeNewlines(md)

	// create markdown parser
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// parse markdown into AST tree
	doc := p.Parse(md)

	// optional: see AST tree
	if false {
		fmt.Printf("%s", "--- AST tree:\n")
		ast.Print(os.Stdout, doc)
	}

	// create HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := markdown.Render(doc, renderer)

	fmt.Printf("\n--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, html)
}
