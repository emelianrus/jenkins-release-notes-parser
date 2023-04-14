package sources

import "github.com/emelianrus/jenkins-release-notes-parser/types"

// Unused for now
type Downloader interface {
	Download(projectName string) ([]types.ReleaseNote, error)
}

// github := &Github{}
// pluginSite := &PluginSiteDownloader{}

// DownloadPlugin(github, "my-github-project")
// DownloadPlugin(pluginSite, "my-plugin-site-project")
// Download single plugin from source
func DownloadPlugin(d Downloader, projectName string) ([]types.ReleaseNote, error) {
	releaseNotes, err := d.Download(projectName)
	return releaseNotes, err
}
