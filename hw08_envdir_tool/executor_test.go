package main

import "testing"

func TestRunCmd(t *testing.T) {
	code := RunCmd([]string{"true"}, Environment{})
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
}
