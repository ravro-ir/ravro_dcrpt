package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func ReadCurrentDir() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	lstFileDecrypt, _ := WalkMatch(path, "*.ravro")
	for _, name := range lstFileDecrypt {
		filename := filepath.Base(name)
		filename = strings.Replace(filename, ".ravro", "", 1)
		_, err := SslDecrypt(name, filename)
		if err != nil {
			fmt.Println("We have error in for decode")
			os.Exit(0)
		}
	}
}
