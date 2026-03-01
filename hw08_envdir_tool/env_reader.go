package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.Contains(name, "=") {
			return nil, errors.New("invalid env var name: contains '='")
		}

		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		// Remove variable first in any case
		ev := EnvValue{NeedRemove: true}

		if len(data) == 0 {
			env[name] = ev
			continue
		}

		// Take first line
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			data = data[:i]
		}

		// Replace terminal zeroes with newline
		data = bytes.ReplaceAll(data, []byte{0x00}, []byte("\n"))

		// Trim trailing spaces and tabs
		value := strings.TrimRight(string(data), " \t")

		ev.Value = value
		env[name] = ev
	}

	return env, nil
}
