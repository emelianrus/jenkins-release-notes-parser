package pluginManager

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	// TODO: enable for windows only
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func TestPluginManager_FixPluginDependencies(t *testing.T) {
	havepm := NewPluginManager()
	havepm.Plugins = map[string]*Plugin{
		"workflow-step-api": NewPluginWithVersion("workflow-step-api", "625.vd896b_f445a_f8"),
		"scm-api":           NewPluginWithVersion("scm-api", "602.v6a_81757a_31d2"),
		"structs":           NewPluginWithVersion("structs", "1.22"),
	}

	wantpm := NewPluginManager()
	wantpm.Plugins = map[string]*Plugin{
		"workflow-step-api": {
			Name:    "workflow-step-api",
			Version: "625.vd896b_f445a_f8",
			Url:     "https://updates.jenkins.io/download/plugins/workflow-step-api/625.vd896b_f445a_f8/workflow-step-api.hpi",
			Type:    UNKNOWN,
			Dependencies: map[string]Plugin{
				"structs": *NewPluginWithVersion("structs", "308.v852b473a2b8c"),
			},
			RequiredBy: make(map[string]string),
		},
		"scm-api": {
			Name:    "scm-api",
			Version: "602.v6a_81757a_31d2",
			Url:     "https://updates.jenkins.io/download/plugins/scm-api/602.v6a_81757a_31d2/scm-api.hpi",
			Type:    UNKNOWN,
			Dependencies: map[string]Plugin{
				"structs": *NewPluginWithVersion("structs", "308.v852b473a2b8c"),
			},
			RequiredBy: make(map[string]string),
		},
		"structs": {
			Name:         "structs",
			Version:      "308.v852b473a2b8c",
			Url:          "https://updates.jenkins.io/download/plugins/structs/308.v852b473a2b8c/structs.hpi",
			Type:         UNKNOWN,
			Dependencies: make(map[string]Plugin),
			RequiredBy: map[string]string{
				"scm-api":           "602.v6a_81757a_31d2",
				"workflow-step-api": "625.vd896b_f445a_f8",
			},
		},
	}

	havepm.FixPluginDependencies()

	diff := cmp.Diff(&havepm.Plugins, &wantpm.Plugins)
	if diff != "" {
		t.Errorf("Not expected result")
		fmt.Println(diff)
	}
}

func TestPluginManager_LoadWarnings(t *testing.T) {
	pm := NewPluginManager()
	pm.Plugins = map[string]*Plugin{
		"blueocean":              NewPluginWithVersion("blueocean", "1.23.2"),
		"hashicorp-vault-plugin": NewPluginWithVersion("hashicorp-vault-plugin", "3.8.0"),
	}

	tests := []struct {
		p    map[string]*Plugin
		want map[string]map[string]string
	}{
		{
			p: pm.Plugins,
			want: map[string]map[string]string{
				"blueocean": {
					"Path traversal vulnerability":                     "1.23.2",
					"Missing permission check":                         "1.23.2",
					"CSRF vulnerability and missing permission checks": "1.25.3",
				},
				"hashicorp-vault-plugin": {
					"Agent-to-controller security bypass":                         "3.8.0",
					"Path traversal vulnerability allows reading arbitrary files": "336.v182c0fbaaeb7",
					"Missing permission checks allow capturing credentials":       "354.vdb_858fd6b_f48",
				},
			},
		},
	}

	for _, tt := range tests {

		for _, plugin := range tt.p {
			plugin.LoadWarnings()
			for _, warn := range plugin.Warnings {
				if _, found := tt.want[plugin.Name][warn.Message]; found {
					if warn.Versions[0].LastVersion != tt.want[plugin.Name][warn.Message] {
						t.Errorf("'%s' Warning version should be '%v'", warn.Versions[0].LastVersion, tt.want[plugin.Name][warn.Message])
					}
				} else {
					t.Errorf("'%s' should not be here we have '%v'", warn.Message, tt.want)
				}
			}

			if len(plugin.Warnings) != len(tt.want[plugin.Name]) {
				t.Errorf("Warnings len '%d' != want len '%d'", len(tt.p[plugin.Name].Warnings), len(tt.want[plugin.Name]))
			}
		}
	}
}

// func TestPluginManager_FixWarnings(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		p    *PluginManager
// 	}{
// 		{
// 			p: &PluginManager{
// 				"blueocean": {
// 					Name:    "blueocean",
// 					Version: "1.23.2",

// 					Url:          "",
// 					Type:         UNKNOWN,
// 					Dependencies: make(map[string]Plugin),
// 					RequiredBy:   make(map[string]string),
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.p.FixWarnings()
// 		})
// 	}
// }

// func TestPluginManager_GetMainPlugins(t *testing.T) {
// 	have := &PluginManager{
// 		"workflow-step-api": {
// 			Name:         "workflow-step-api",
// 			Version:      "625.vd896b_f445a_f8",
// 			Url:          "",
// 			Type:         UNKNOWN,
// 			Dependencies: make(map[string]Plugin),
// 			RequiredBy:   make(map[string]string),
// 		},
// 		"scm-api": {
// 			Name:         "scm-api",
// 			Version:      "602.v6a_81757a_31d2",
// 			Url:          "",
// 			Type:         UNKNOWN,
// 			Dependencies: make(map[string]Plugin),
// 			RequiredBy:   make(map[string]string),
// 		},
// 		"structs": {
// 			Name:         "structs",
// 			Version:      "1.22",
// 			Url:          "",
// 			Type:         UNKNOWN,
// 			Dependencies: make(map[string]Plugin),
// 			RequiredBy:   make(map[string]string),
// 		},
// 	}

// 	//want := []string{"workflow-step-api", "scm-api"}
// 	have.FixPluginDependencies()

// 	have.GetStandalonePlugins()

// }

func TestNewPluginManager(t *testing.T) {
	tests := []struct {
		name string
		want PluginManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPluginManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPluginManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginManager_AddPlugin(t *testing.T) {
	pm := NewPluginManager()

	tests := []struct {
		name string
		pm   *PluginManager
		args []Plugin
	}{
		{
			name: "test add plugins",
			pm:   &pm,
			args: []Plugin{
				*NewPluginWithVersion("blueocean", "1.23.3"),
				*NewPluginWithVersion("configuration-as-code", "1616.v11393eccf675"),
				*NewPluginWithVersion("hashicorp-vault-plugin", "3.8.0"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, pluginName := range tt.args {
				tt.pm.AddPlugin(&pluginName)
			}
		})
	}
}
