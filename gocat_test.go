package gocat

import (
	"strings"
	"testing"
)

func TestCopyStream(t *testing.T) {
	var tests = []struct {
		input string
	}{
		{"abc"},
	}
	for _, test := range tests {
		src := strings.NewReader(test.input)
		dst := new(strings.Builder)
		CopyStream(src, dst)

		if test.input != dst.String() {
			t.Errorf("CopyStream(%s, _) => %s", test.input, dst.String())
		}
	}
}
