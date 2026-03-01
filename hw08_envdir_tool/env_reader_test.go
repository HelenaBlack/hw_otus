package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "FOO"), []byte("bar \t\nbaz"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := env["FOO"]
	if !ok {
		t.Fatalf("FOO not found")
	}

	if val.Value != "bar" {
		t.Fatalf("unexpected value: %q", val.Value)
	}
}
