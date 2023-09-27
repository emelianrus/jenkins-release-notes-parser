package parsers

import (
	"strings"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/sirupsen/logrus"
)

type InputParser interface {
	Parse(content []byte) map[string]string
}
type TXTFileParser struct{}

// parse plugins.txt with types
func (fp TXTFileParser) Parse(content []byte) map[string]string {

	// pluginManager := NewPluginManager()
	resultMap := make(map[string]string)

	for _, line := range strings.Split(string(content), "\n") {
		var name, version string

		switch {
		// skip commented and empty lines
		case strings.HasPrefix(line, "#") || line == "":
			continue

		// TODO:
		// case strings.Contains(line, "::"):
		// 	parsedLine := strings.Split(line, "::")
		// 	name = parsedLine[0]

		// 	if utils.IsUrl(parsedLine[1]) {
		// 		url = parsedLine[1]
		// 		pluginManager.AddPlugin(NewPluginWithUrl(name, url))
		// 	} else {
		// 		logrus.Errorln("Unsupported format NAME::VERSION use instead NAME:VERSION")
		// 	}
		// 	continue
		case strings.Contains(line, ":"):
			parsedLine := strings.Split(line, ":")
			name = parsedLine[0]

			if utils.IsUrl(parsedLine[1]) {
				logrus.Errorln("Unsupported format NAME:URL use instead NAME::URL")

			} else {
				version = strings.TrimRight(parsedLine[1], "\r")
				resultMap[name] = version
			}
			continue

		default:
			logrus.Warnf("Line '%s' is unsupported, skipped.\n", line)
		}
	}

	return resultMap

}
