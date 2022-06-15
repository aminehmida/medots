package shell

import (
	"bytes"
	"os"
	"os/exec"
)

func Run(command string, shellToUse string, interactive bool) (string, string, error) {
	if interactive {
		cmd := exec.Command(shellToUse, "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		return "", "", err
	} else {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command(shellToUse, "-c", command)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		return stdout.String(), stderr.String(), err
	}
}
