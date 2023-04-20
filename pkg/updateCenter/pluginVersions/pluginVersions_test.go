package pluginVersions

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	gotd, _ := Get()
	fmt.Println(gotd.Plugins["blueocean"]["1.23.3"].RequiredCore)

	tests := []struct {
		name    string
		want    *PluginVersions
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get()
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
