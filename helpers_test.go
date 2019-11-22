package main

import (
	"testing"
	"strings"
)

func TestCopyStream(t *testing.T) {
	var tests = []struct{
		input string
		chunk_size int
	} {
		{"abc", 1},
		{"abc", 2},
		{"abc", 3},
		{"abc", 4},
	}
	for _, test := range tests {
		src := strings.NewReader(test.input)
		dst := new(strings.Builder)
		CopyStream(src, dst, test.chunk_size)

		if test.input != dst.String() {
			t.Errorf("CopyStream(%s, _, %d) => %s", test.input, test.chunk_size, dst.String())
		}
	}

}
