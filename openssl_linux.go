package main

import (
	"bytes"
	"os/exec"
)

func Ssl_decrypt(command string) (out string, err_out error) {

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	name := "bash"
	arg := []string{"-c", command}
	cmd := exec.Command(name, arg...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return cmd.String(), nil
}
