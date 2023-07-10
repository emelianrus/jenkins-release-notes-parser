package utils

import (
	"testing"
)

func TestIsNewerThan(t *testing.T) {

	var tests = []struct {
		ver1, ver2 string
		want       bool
	}{
		{"2.1.1", "1.1.1", true},
		{"1.1.1.1", "1.1.1", true},
		{"1.1.1.4", "1.1.1-5.0", true},
		{"1.1.1-2.0", "1.1.1-1.0", true},
		{"4.5.13-1.0", "4.5.10-2.0", true},
		{"1108.v57edf648f5d4", "1.1.1", true},
		{"1208.v57edf648f5d4", "1108.v57edf648f5d4", true},
		{"2.13.2.20220328-281.v9ecc7a_5e834f", "2.13.2.20220328-273.v11d70a_b_a_1a_52", true},
		{"5.11.2-182.v0f1cf4c5904e", "5.4.2", true},
		{"5.4.1-beta-1", "4.13.2-1", true},

		{"1.1.1", "1.1.1", false},
		{"1.1.1", "2.1.1", false},
		{"1.1.1", "1.1.1.1", false},
		{"1.1.1-5.0", "1.1.1.4", false},
		{"1.1.1-1.0", "1.1.1-2.0", false},
		{"4.5.10-2.0", "4.5.13-1.0", false},
		{"1.1.1", "1108.v57edf648f5d4", false},
		{"1108.v57edf648f5d4", "1108.v57edf648f5d4", false},
		{"1108.v57edf648f5d4", "1208.v57edf648f5d4", false},
		{"2.13.2.20220328-273.v11d70a_b_a_1a_52", "2.13.2.20220328-281.v9ecc7a_5e834f", false},
		{"5.4.2", "5.11.2-182.v0f1cf4c5904e", false},
		{"4.13.2-1", "5.4.1-beta-1", false},
	}

	for _, test := range tests {
		if output := IsNewerThan(test.ver1, test.ver2); output != test.want {
			t.Errorf("Output '%t' not equal to want '%t' for ver1: %s and ver2: %s", output, test.want, test.ver1, test.ver2)
		}
	}
}
