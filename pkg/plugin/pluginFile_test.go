package plugin

import (
	"reflect"
	"testing"
)

func TestParsePlugins(t *testing.T) {
	tests := []struct {
		name string
		f    *File
		want *PluginManager
	}{
		{
			name: "Test plugin with Version",
			f: &File{
				Name:    "testfile",
				Content: []byte("blueocean:1.24.2"),
			},
			want: &PluginManager{
				"blueocean": NewPluginWithVersion("blueocean", "1.24.2"),
			},
		},
		{
			name: "Test plugin with URL",
			f: &File{
				Name:    "testfile",
				Content: []byte("blueocean::https://updates.jenkins.io/download/plugins/blueocean/1.23.2/blueocean.hpi"),
			},
			want: &PluginManager{
				"blueocean": NewPluginWithUrl("blueocean", "https://updates.jenkins.io/download/plugins/blueocean/1.23.2/blueocean.hpi"),
			},
		},
		{
			name: "Test two plugins in list",
			f: &File{
				Name:    "testfile",
				Content: []byte("blueocean:1.24.2\nworkflow-step-api:2.22"),
			},
			want: &PluginManager{
				"blueocean":         NewPluginWithVersion("blueocean", "1.24.2"),
				"workflow-step-api": NewPluginWithVersion("workflow-step-api", "2.22"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePlugins(tt.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePlugins() = %v, want %v", got, tt.want)
			}
		})
	}
}
