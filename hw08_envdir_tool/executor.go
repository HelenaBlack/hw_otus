package main

import (
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Start with current environment
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

	// Apply envdir rules
	for key, val := range env {
		delete(envMap, key)
		if val.Value != "" {
			envMap[key] = val.Value
		}
	}

	// Convert back to slice
	finalEnv := make([]string, 0, len(envMap))
	for k, v := range envMap {
		finalEnv = append(finalEnv, k+"="+v)
	}

	command.Env = finalEnv

	err := command.Run()
	if err == nil {
		return 0
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}

	return 1
}
