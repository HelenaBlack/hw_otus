package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 111
	}

	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		if i := len(e); i > 0 {
			parts := []rune(e)
			for j := 0; j < len(parts); j++ {
				if parts[j] == '=' {
					envMap[string(parts[:j])] = string(parts[j+1:])
					break
				}
			}
		}
	}

	for key, val := range env {
		delete(envMap, key)
		if val.Value != "" {
			envMap[key] = val.Value
		}
	}

	finalEnv := make([]string, 0, len(envMap))
	for k, v := range envMap {
		finalEnv = append(finalEnv, k+"="+v)
	}

	command.Env = finalEnv

	err := command.Run()
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}

	return 1
}
