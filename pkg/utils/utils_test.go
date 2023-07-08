package utils

import (
	"testing"
)

func TestIsUrl(t *testing.T) {
	var tests = []struct {
		url  string
		want bool
	}{
		{"https://urlwithhttps.com", true},
		{"http://urlwithhttp.com", true},
		{"notaurl.forsure", false},
	}

	for _, test := range tests {
		if output := IsUrl(test.url); output != test.want {
			t.Errorf("Output '%t' not equal to want '%t' for url: %s", output, test.want, test.url)
		}
	}
}

func TestSortVersions(t *testing.T) {

	versions := []string{
		"1.19.1",
		"1.15.0",
		"1.17.0",
		"1.24.7",
		"1.10.1",
		"1.25.3",
		"1.24.6",
		"1.25.0",
		"1.24.2",
		"1.26.0",
		"1.25.1",
		"1.24.1",
		"1.25.5",
		"1.25.2",
		"1.25.6",
		"1.24.8",
		"1.19.2",
		"1.23.1",
		"1.25.4",
		"1.25.7",
		"1.25.8",
		"1.24.0",
		"1.22.0",
		"1.24.3",
		"1.15.1",
		"1.18.1",
		"1.23.3",
		"1.24.4",
		"1.24.5",
		"1.23.0",
		"1.13.0",
		"1.14.0",
		"1.19.0",
		"1.23.2",
		"1.21.0",
	}

	var tests = []struct {
		want    string
		have    string
		IsError bool
	}{
		{
			have:    "1.25.8",
			want:    "1.26.0",
			IsError: false,
		},
		{
			have:    "1.26.0",
			want:    "",
			IsError: true,
		},
		{
			have:    "notexist",
			want:    "",
			IsError: true,
		},
	}

	for _, test := range tests {
		nextVersion, err := GetNextVersion(versions, test.have)

		if nextVersion != test.want {
			t.Errorf("Version have %s want %s err: %v", nextVersion, test.want, err)
		}
		// TODO: is there another way?
		// dont want to do mess with bools
		foundError := false

		if err == nil {
			foundError = false
		} else {
			foundError = true
		}

		if test.IsError != foundError {
			t.Errorf("Error have %t want %t ", test.IsError, foundError)
		}
	}
}
