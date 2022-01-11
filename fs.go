package main

import (
	"errors"
	"fmt"
	"golang.org/x/exp/utf8string"
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
		newName := strings.ReplaceAll(name, " ", "")
		err := os.Rename(name, newName)
		if err != nil {
			log.Fatalf("We have error in %s", err)
		}
		name = newName
		filename := filepath.Base(name)
		asciiCheck := utf8string.NewString(filename).IsASCII()
		if !asciiCheck {
			fmt.Println(asciiCheck)
		}
		fmt.Println(asciiCheck)
		filename = strings.Replace(filename, ".ravro", "", 1)
		getFileName := fileNameWithoutExtension(filename)
		fileExt := filepath.Ext(filename)
		fmt.Println(getFileName)
		fmt.Println(fileExt)
		_, err = SslDecrypt(name, filename)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
