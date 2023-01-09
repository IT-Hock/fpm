package utils

import (
	"bytes"
	"os/exec"
)

func RunCommand(command string, args ...string) (*exec.Cmd, string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return cmd, "", err
	}

	return cmd, out.String(), nil
}
