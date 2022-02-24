package releaseParser

import (
	"strings"

	types "github.com/emelianrus/jenkins-release-notes-parser/types"
)

func GetChangedVersions(diff *types.Diff) (*types.Project, error) {
	project := &types.Project{}

	for _, f := range diff.Files {
		file := &types.File{}
		file.Name = f.OrigName
		project.Files = append(project.Files, file)
		elementMap := make(map[string]types.Plugin)

		for _, h := range f.Chunks {
			plugin := types.Plugin{}
			var itemExist bool

			for _, l := range h.Lines {
				if l.Mode == types.UNCHANGED {
					continue
				}
				pluginSplited := strings.Split(l.Content, ":")
				// unsupported string
				if len(pluginSplited) != 2 {
					continue
				}

				if _, found := elementMap[pluginSplited[0]]; found {
					plugin = elementMap[pluginSplited[0]]
					itemExist = true
				} else {
					itemExist = false
				}

				plugin.Name = pluginSplited[0]

				if !itemExist {
					plugin.OldVersion = pluginSplited[1]
				}

				if itemExist {
					plugin.NewVersion = pluginSplited[1]
				}

				elementMap[pluginSplited[0]] = plugin
			}
		}

		for _, v := range elementMap {
			if v.NewVersion != "" && v.OldVersion != "" {
				file.Plugins = append(file.Plugins, v)
			}
		}
	}

	return project, nil
}
