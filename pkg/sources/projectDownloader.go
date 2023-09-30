package sources

import (
	"sort"

	"github.com/emelianrus/jenkins-release-notes-parser/pkg/utils"
	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/sirupsen/logrus"
)

// Unused for now
type ReleaseNotesDownloader interface {
	Download(projectName string) ([]types.ReleaseNote, error)
}

// github := &Github{}
// pluginSite := &PluginSiteDownloader{}

// DownloadPlugin(github, "my-github-project")
// DownloadPlugin(pluginSite, "my-plugin-site-project")
// Download single plugin from source
func DownloadProjectReleaseNotes(d ReleaseNotesDownloader, projectName string) ([]types.ReleaseNote, error) {
	logrus.Infoln("[DownloadProjectReleaseNotes] started with")
	releaseNotes, err := d.Download(projectName)

	sort.Slice(releaseNotes, func(i, j int) bool {
		return utils.IsNewerThan(releaseNotes[i].Name, releaseNotes[j].Name)
	})

	return releaseNotes, err
}
