package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func SslDecrypt(name, filename string) (out string, errOut error) {

	command := fmt.Sprintf("openssl smime -decrypt -in %s -inform DER -inkey key/key.private -out decrypt/%s -binary", name, filename)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	name = "bash"
	arg := []string{"-c", command}
	cmd := exec.Command(name, arg...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	return cmd.String(), nil
}
