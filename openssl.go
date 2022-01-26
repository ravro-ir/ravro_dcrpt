package main

import (
	"fmt"
	"os/exec"
)

func SslDecrypt(name, filename string) (out string, errOut error) {
	args := []string{"smime", "-decrypt", "-in", name, "-inform", "DER", "-inkey", "key/key.private", "-out", "decrypt/" + filename, "-binary"}
	output, err_ := RunCMD("openssl", args, true)
	if err_ != nil {
		return "", err_
	}
	return output, nil
}

func RunCMD(path string, args []string, debug bool) (out string, err error) {

	cmd := exec.Command(path, args...)
	var b []byte
	b, err = cmd.CombinedOutput()
	out = string(b)
	if err != nil {
		return "", fmt.Errorf("The private key is incorrect")
	}
	return out, nil
}
