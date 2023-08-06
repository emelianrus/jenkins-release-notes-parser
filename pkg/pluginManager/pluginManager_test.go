package pluginManager

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

func init() {
	logLevel := os.Getenv("RN_DEBUG")
	if logLevel != "" {
		lvl, _ := logrus.ParseLevel(logLevel)
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// TODO: enable for windows only
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
}

func TestPluginManager_FixPluginDependencies(t *testing.T) {
	havepm := NewPluginManager()
	havepm.AddPluginWithVersion("workflow-step-api", "625.vd896b_f445a_f8")
	havepm.AddPluginWithVersion("scm-api", "602.v6a_81757a_31d2")
	havepm.AddPluginWithVersion("structs", "1.22")

	wantpm := NewPluginManager()
	wantpm.Plugins = map[string]*Plugin{
		"workflow-step-api": {
			Name:                "workflow-step-api",
			Version:             "625.vd896b_f445a_f8",
			RequiredCoreVersion: "2.289.1",
			Url:                 "https://updates.jenkins.io/download/plugins/workflow-step-api/625.vd896b_f445a_f8/workflow-step-api.hpi",
			Dependencies:        map[string]Plugin{"structs": {Name: "structs", Version: "308.v852b473a2b8c"}},
			RequiredBy:          make(map[string]string),
		},
		"scm-api": {
			Name:                "scm-api",
			Version:             "602.v6a_81757a_31d2",
			RequiredCoreVersion: "2.289.1",
			Url:                 "https://updates.jenkins.io/download/plugins/scm-api/602.v6a_81757a_31d2/scm-api.hpi",
			Dependencies:        map[string]Plugin{"structs": {Name: "structs", Version: "308.v852b473a2b8c"}},
			RequiredBy:          make(map[string]string),
		},
		"structs": {
			Name:                "structs",
			Version:             "308.v852b473a2b8c",
			RequiredCoreVersion: "2.222.4",
			Url:                 "https://updates.jenkins.io/download/plugins/structs/308.v852b473a2b8c/structs.hpi",
			Dependencies:        make(map[string]Plugin),
			RequiredBy:          map[string]string{},
		},
	}

	havepm.FixPluginDependencies()
	for _, v := range havepm.UpdatedPlugins {
		v.LatestVersion = ""
		v.GITUrl = ""
	}
	diff := cmp.Diff(&wantpm.Plugins, &havepm.UpdatedPlugins)
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
					"Improper masking of credentials":                             "360.v0a_1c04cf807d",
				},
			},
		},
	}

	pm.LoadWarnings()

	for _, tt := range tests {
		for _, plugin := range tt.p {
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

func TestPluginManager_FixWarnings(t *testing.T) {
	pm := NewPluginManager()

	pm.AddPluginWithVersion("blueocean", "1.23.2")
	tests := []struct {
		name string
		p    *PluginManager
	}{
		{
			p: &pm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.FixWarnings()

			fmt.Println(tt.p.Plugins["blueocean"].Version)
		})
	}
}

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
		args map[string]string
	}{
		{
			name: "test add plugins",
			pm:   &pm,
			args: map[string]string{
				"blueocean":              "1.23.3",
				"configuration-as-code":  "1616.v11393eccf675",
				"hashicorp-vault-plugin": "3.8.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for name, ver := range tt.args {
				tt.pm.AddPluginWithVersion(name, ver)
			}
		})
	}
}

func TestPluginManager_SetCoreVersion(t *testing.T) {
	pm := NewPluginManager()
	newVer := "1.1.1"
	pm.SetCoreVersion(newVer)

	t.Run("test set new version", func(t *testing.T) {
		ver := pm.GetCoreVersion()
		if ver != newVer {
			t.Errorf("SetCoreVersion() = %v, want %v", ver, newVer)
		}
	})
}

func TestPluginManager_AddPluginWithVersion(t *testing.T) {
	type args struct {
		pluginName string
		version    string
	}
	pm := NewPluginManager()

	tests := []struct {
		name string
		pm   *PluginManager
		args args
	}{
		{
			name: "check add new plugin",
			pm:   &pm,
			args: args{
				pluginName: "blueocean",
				version:    "1.23.2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm.AddPluginWithVersion(tt.args.pluginName, tt.args.version)
			plugins := pm.GetPlugins()

			if _, found := plugins[tt.args.pluginName]; !found {
				t.Error("added plugin is not in plugin manager list")
			}

			if plugins[tt.args.pluginName].Version != tt.args.version {
				t.Error("added plugin version is not correct")
			}

		})
	}
}
