package main

import (
	"testing"

	"golang.org/x/example/hello/reverse"
)

func TestReverseString(t *testing.T) {
	result := reverse.String("Hello, OTUS!")
	expected := "!SUTO ,olleH"
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
