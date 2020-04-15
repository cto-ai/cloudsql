package osutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecCmd executes shell command
func ExecCmd(command string) error {
	if len(command) < 1 {
		return fmt.Errorf("Error: command is empty")
	}

	splitCmd := strings.Split(command, " ")
	bin := splitCmd[0]
	args := splitCmd[1:]

	cmd := exec.Command(bin, args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error executing command `%v` %v", command, err)
	}
	return nil
}

// ExecCmdWithLogs executes shell command returning stdout
func ExecCmdWithLogs(command string) (string, error) {
	if len(command) < 1 {
		return "", fmt.Errorf("Error: command is empty")
	}

	splitCmd := strings.Split(command, " ")
	bin := splitCmd[0]
	args := splitCmd[1:]

	cmd := exec.Command(bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("Error executing command `%v`: %v", command, stderr.String())
	}
	// fmt.Println("out:", stdout.String(), "err:", stderr.String())

	return stdout.String(), nil
}
