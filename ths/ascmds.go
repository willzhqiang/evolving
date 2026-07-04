package ths

import (
	"fmt"
	"os/exec"
	"strings"
)

func Run(script string) (string, error) {
	cmd := commandForScript(script)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return strings.TrimSpace(string(exitErr.Stderr)), fmt.Errorf("osascript error: %w", exitErr)
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func commandForScript(script string) *exec.Cmd {
	trimmed := strings.TrimSpace(script)
	if strings.HasPrefix(trimmed, "osascript ") {
		return exec.Command("/bin/zsh", "-lc", trimmed)
	}
	return exec.Command("/usr/bin/osascript", "-e", script)
}
