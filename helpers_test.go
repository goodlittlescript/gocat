package main

import (
	"testing"
	"strings"
)

func TestCopyStream(t *testing.T) {
	actual := "Lorem ipsum"
	src := strings.NewReader(actual)
	dst := new(strings.Builder)
	CopyStream(src, dst, 1)
	expected := dst.String()
	if actual != expected {
		t.Errorf("Failed %s != %s", actual, expected)
	}
}
