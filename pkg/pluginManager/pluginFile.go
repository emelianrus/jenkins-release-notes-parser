package pluginManager

// this file should only be for parse file
/*

TODO: add conventions for package parser

like

:::
vs ::
vs pluginname:version
vs pluginname:url
else ...

TODO:

should be able read plugins from yaml file
should be able read plugins from cli arg

*/

import (
	"fmt"
	"os"
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/sirupsen/logrus"
)

type File struct {
	Name    string
	Content []byte
}

// TODO: remove redundant
func ReadFile(fileName string) *File {
	// Read file content
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	return &File{
		Name:    fileName,
		Content: fileContent,
	}
}

// parse plugins.txt with types
func ParsePlugins(f *File) *PluginManager {

	pluginManager := NewPluginManager()

	for _, line := range strings.Split(string(f.Content), "\n") {
		var name, version, url string

		switch {
		// skip commented and empty lines
		case strings.HasPrefix(line, "#") || line == "":
			continue

		// TODO: DRY
		case strings.Contains(line, "::"):
			parsedLine := strings.Split(line, "::")
			name = parsedLine[0]

			if utils.IsUrl(parsedLine[1]) {
				url = parsedLine[1]
				pluginManager.AddPlugin(NewPluginWithUrl(name, url))
			} else {
				logrus.Errorln("Unsupported format NAME::VERSION use instead NAME:VERSION")
			}
			continue
		case strings.Contains(line, ":"):
			parsedLine := strings.Split(line, ":")
			name = parsedLine[0]

			if utils.IsUrl(parsedLine[1]) {
				logrus.Errorln("Unsupported format NAME:URL use instead NAME::URL")

			} else {
				version = parsedLine[1]
				pluginManager[name] = NewPluginWithVersion(name, version)
			}
			continue

		default:
			logrus.Warnf("Line '%s' is unsupported, skipped.\n", line)
		}
	}

	return &pluginManager

}

/* TODO:

add func

* replace lines at file to update pluginFile
how to append?

*/

// TODO: REFACT
// find and replace line in string, return true if changed false if unchanged
func findReplace(name string, version string, lines []string) ([]string, bool) {
	for i, line := range lines {
		if strings.HasPrefix(line, name+":") {
			lines[i] = name + ":" + version
			return lines, true
		}
	}
	return lines, false
}

// TODO: CHECK || REFACT
func (p *PluginManager) UpdateFile(fileName string) {
	fmt.Println("Update current file: " + fileName + ".updated")

	lines := strings.Split(string(ReadFile(fileName).Content), "\n")

	plugins := make(map[string]Plugin)
	for _, v := range *p {
		plugins[v.Name] = *v
	}

	transition := make(map[string]Plugin)
	for _, plugin := range plugins {
		linesChanged, replaced := findReplace(plugin.Name, plugin.Version, lines)
		if replaced {
			lines = linesChanged
		} else {
			transition[plugin.Name] = plugin
		}
	}

	output := strings.Join(lines, "\n")

	transitionPlugins := ""
	for _, plugin := range transition {
		transitionPlugins += plugin.Name + ":" + plugin.Version + "\n"
	}

	output += transitionPlugins

	err := os.WriteFile(fileName+".updated", []byte(output), 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
}

// TODO: CHECK REFACT
// func (p *PluginManager) WriteToFile(outputFile string) {
// 	if p.FileName == outputFile {
// 		p.UpdateFile(outputFile)
// 		return
// 	}

// 	file, err := os.Create(outputFile)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer file.Close()

// 	w := bufio.NewWriter(file)

// 	w.WriteString("# === EXTERNAL ===\n")
// 	for _, pluginL := range p.External {
// 		w.WriteString(pluginL.Name + ":" + pluginL.Version + "\n")
// 	}

// 	w.WriteString("\n# === PRIMARY ===\n")

// 	keysP := make([]string, 0, len(p.External))
// 	for k := range p.Primary {
// 		keysP = append(keysP, k)
// 	}
// 	sort.Strings(keysP)

// 	for _, k := range keysP {
// 		w.WriteString(p.Primary[k].Name + ":" + p.Primary[k].Version + "\n")
// 	}

// 	w.WriteString("\n# === TRANSITIVE ===\n")

// 	keys := make([]string, 0, len(p.External))
// 	for k := range p.Transitive {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	for _, k := range keys {
// 		w.WriteString(p.Transitive[k].Name + ":" + p.Transitive[k].Version + "\n")
// 	}

// 	w.Flush()
// }
