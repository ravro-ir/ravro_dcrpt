package main

import (
	"errors"
	"fmt"
	"golang.org/x/exp/utf8string"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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
		var NewPathFile string
		var oldNamePath string
		var newNamePath string
		newName := strings.ReplaceAll(name, " ", "")
		err := os.Rename(name, newName)
		if err != nil {
			log.Fatalf("We have error in %s", err)
		}
		name = newName
		filename := filepath.Base(name)
		basePath := filepath.Dir(name)
		filename = strings.Replace(filename, ".ravro", "", 1)
		fileExt := filepath.Ext(filename)
		oldName := fileNameWithoutExtension(filename)
		asciiCheck := utf8string.NewString(filename).IsASCII()
		if !asciiCheck {
			newNameRand := randSeq(10)
			if runtime.GOOS == "windows" {
				NewPathFile = basePath + "\\" + newNameRand + fileExt + ".ravro"
			} else {
				NewPathFile = basePath + "/" + newNameRand + fileExt + ".ravro"
			}
			err := os.Rename(name, NewPathFile)
			if err != nil {
				return
			}
			filename = newNameRand + fileExt
			name = NewPathFile
		}
		_, err = SslDecrypt(name, filename)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		if runtime.GOOS == "windows" {
			oldNamePath = "datadecrypt" + "\\" + oldName + fileExt
			newNamePath = "datadecrypt" + "\\" + filename
		} else {
			oldNamePath = "datadecrypt" + "/" + oldName + fileExt
			newNamePath = "datadecrypt" + "/" + filename
		}
		err = os.Remove(name)
		if err != nil {
			return
		}
		err = os.Rename(newNamePath, oldNamePath)
		if err != nil {
			return
		}
	}
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
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
