package pluginManager

import (
	"fmt"
	"os"
	"reflect"
	"testing"

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

func TestPlugin_Download(t *testing.T) {
	plugin := &Plugin{
		Name:    "blueocean",
		Version: "1.23.3",
	}
	tests := []struct {
		name    string
		p       *Plugin
		want    string
		wantErr bool
	}{
		{
			name: "blueocean",
			p:    plugin,
			want: "plugins/blueocean-1.23.3.hpi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Download()
			if (err != nil) != tt.wantErr {
				t.Errorf("Plugin.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Plugin.Download() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlugin_LoadDependenciesFromManifest(t *testing.T) {
	pl := NewPluginWithVersion("scm-api", "602.v6a_81757a_31d2")
	tests := []struct {
		name string
		p    *Plugin
		want map[string]string
	}{
		{
			name: "scm-api",
			p:    pl,
			want: map[string]string{
				"structs": "308.v852b473a2b8c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.LoadDependenciesFromManifest()

			fmt.Println(tt.p.Dependencies)

			for _, dep := range tt.p.Dependencies {
				if _, found := tt.want[dep.Name]; !found {
					t.Errorf("%s should not be here we have %v", dep.Name, tt.want)
				} else {
					if dep.Version != tt.want[dep.Name] {
						t.Errorf("%s Dependenc version should be %v", dep.Version, tt.want[dep.Name])
					}
				}
			}

			if len(tt.p.Dependencies) != len(tt.want) {
				t.Errorf("Dependencies len %d != want len %d", len(tt.p.Warnings), len(tt.want))
			}
		})
	}
}

func TestPlugin_LoadRequiredCoreVersion(t *testing.T) {
	tests := []struct {
		name string
		p    *Plugin
		want string
	}{
		{
			name: "blueocean",
			p: &Plugin{
				Name:    "blueocean",
				Version: "1.23.2",
			},
			want: "2.150.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.p.LoadRequiredCoreVersion()
			if tt.p.RequiredCoreVersion != tt.want {
				t.Errorf("Want version %s differs to have %s", tt.want, tt.p.RequiredCoreVersion)
			}
		})
	}
}

func TestPlugin_GetManifestAttrs(t *testing.T) {
	tests := []struct {
		name string
		p    *Plugin
		want map[string]string
	}{
		{
			name: "blueocean",
			p: &Plugin{
				Name:    "blueocean",
				Version: "1.23.2",
			},
			want: map[string]string{
				"Url":                     "https://github.com/jenkinsci/blueocean-plugin/blob/master/blueocean/doc/BlueOcean.adoc",
				"Short-Name":              "blueocean",
				"Plugin-ScmUrl":           "https://github.com/jenkinsci/blueocean-plugin",
				"Long-Name":               "Blue Ocean",
				"Specification-Title":     "Blue Ocean is a new project that rethinks the user experience of Jenkins. Designed from the ground up for Jenkins Pipeline and compatible with Freestyle jobs, Blue Ocean reduces clutter and increases clarity for every member of your team.",
				"Hudson-Version":          "2.150.3",
				"Plugin-Version":          "1.23.2",
				"Archiver-Version":        "Plexus Archiver",
				"Build-Jdk":               "1.8.0_231",
				"Implementation-Title":    "blueocean",
				"Implementation-Version":  "1.23.2",
				"Support-Dynamic-Loading": "true",
				"Manifest-Version":        "1.0",
				"Extension-Name":          "blueocean",
				"Group-Id":                "io.jenkins.blueocean",
				"Plugin-Developers":       "Thorsten Iberian Sumurai:scherler:,Cliff Meyers:cliffmeyers:,Tom Fennelly:tfennelly:,Vivek Pandey:vivek:,Kohsuke:kohsuke:,Josh McDonald:sophistifunk:,Ivan Meredith:imeredith:,Michael Neale:michaelneale:,Keith Zantow:kzantow:,James Dumay:i386:,Marc:marcesher:,Paul Dragoonis:dragoonis:,Ivan Santos:pragmaticivan:,Peter Dave Hello:PeterDaveHello:,Alexandru Somai:alexsomai:",
				"Minimum-Java-Version":    "1.8",
				"Jenkins-Version":         "2.150.3",
				"Plugin-License-Url":      "https://opensource.org/licenses/MIT",
				"Built-By":                "bitwiseman",
				"Created-By":              "Apache Maven",
				"Plugin-Dependencies":     "blueocean-bitbucket-pipeline:1.23.2,blueocean-commons:1.23.2,blueocean-config:1.23.2,blueocean-core-js:1.23.2,blueocean-dashboard:1.23.2,blueocean-events:1.23.2,blueocean-git-pipeline:1.23.2,blueocean-github-pipeline:1.23.2,blueocean-i18n:1.23.2,blueocean-jira:1.23.2,blueocean-jwt:1.23.2,blueocean-personalization:1.23.2,blueocean-pipeline-api-impl:1.23.2,blueocean-pipeline-editor:1.23.2,blueocean-rest-impl:1.23.2,blueocean-rest:1.23.2,blueocean-web:1.23.2,jenkins-design-language:1.23.2,blueocean-autofavorite:1.2.3,blueocean-display-url:2.2.0,pipeline-build-step:2.7,pipeline-milestone-step:1.3.1",
				"Plugin-License-Name":     "MIT License",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetManifestAttrs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.GetManifestAttrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
