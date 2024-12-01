package main

import (
	"strings"
	"testing"
)

func TestBuildMermaidGraph(t *testing.T) {
	commits := []Commit{
		{Hash: "a1", Parent: ""},
		{Hash: "b2", Parent: "a1"},
		{Hash: "c3", Parent: "b2"},
	}

	expected := "graph TD\n    a1\n    a1 --> b2\n    b2 --> c3\n"
	result := buildMermaidGraph(commits)

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}
