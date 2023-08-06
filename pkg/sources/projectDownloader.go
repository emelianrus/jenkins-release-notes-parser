package sources

import (
	"fmt"
	"sort"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/sirupsen/logrus"
)

// Unused for now
type Downloader interface {
	Download(projectName string) ([]types.ReleaseNote, error)
}

// github := &Github{}
// pluginSite := &PluginSiteDownloader{}

// DownloadPlugin(github, "my-github-project")
// DownloadPlugin(pluginSite, "my-plugin-site-project")
// Download single plugin from source
func DownloadProjectReleaseNotes(d Downloader, projectName string) ([]types.ReleaseNote, error) {
	logrus.Infoln("[DownloadProjectReleaseNotes] started with")
	releaseNotes, err := d.Download(projectName)
	for _, v := range releaseNotes {
		fmt.Println(v.Name)
	}
	sort.Slice(releaseNotes, func(i, j int) bool {
		return utils.IsNewerThan(releaseNotes[i].Name, releaseNotes[j].Name)
	})

	return releaseNotes, err
}
