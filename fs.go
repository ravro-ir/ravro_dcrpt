package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/utf8string"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ProccesFile struct {
	NewPathFile string
	oldNamePath string
	newNamePath string
	fileExt     string
	oldName     string
	filename    string
	basePath    string
}

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

func AddDir(name string) {
	if err := ensureDir(name); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}
}

func projectpath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return path, err
	}
	return path, nil
}

func fileProccessing(name string) (ProccesFile, error) {

	var processFile ProccesFile

	newName := strings.ReplaceAll(name, " ", "")
	err := os.Rename(name, newName)
	if err != nil {
		return processFile, err
	}
	name = newName
	processFile.filename = filepath.Base(name)
	processFile.basePath = filepath.Dir(name)
	processFile.filename = strings.Replace(processFile.filename, ".ravro", "", 1)
	processFile.fileExt = filepath.Ext(processFile.filename)
	processFile.oldName = fileNameWithoutExtension(processFile.filename)
	asciiCheck := utf8string.NewString(processFile.filename).IsASCII()
	if !asciiCheck {
		newNameRand := randSeq(10)
		if runtime.GOOS == "windows" {
			processFile.NewPathFile = processFile.basePath + "\\" + newNameRand + processFile.fileExt + ".ravro"
		} else {
			processFile.NewPathFile = processFile.basePath + "/" + newNameRand + processFile.fileExt + ".ravro"
		}
		err := os.Rename(name, processFile.NewPathFile)
		if err != nil {
			return processFile, err
		}
		processFile.filename = newNameRand + processFile.fileExt
		name = processFile.NewPathFile
	}
	return processFile, nil
}

func CheckPlatform(process ProccesFile) ProccesFile {
	if runtime.GOOS == "windows" {
		process.oldNamePath = "decrypt" + "\\" + process.oldName + process.fileExt
		process.newNamePath = "decrypt" + "\\" + process.filename
	} else {
		process.oldNamePath = "decrypt" + "/" + process.oldName + process.fileExt
		process.newNamePath = "decrypt" + "/" + process.filename
	}
	return process
}

func JsonParser(process ProccesFile, AnyStruct interface{}) (interface{}, error) {

	jsonFile, err := os.Open(process.oldNamePath)
	reportValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(reportValue, &AnyStruct)
	if err != nil {
		return AnyStruct, err
	}
	err = jsonFile.Close()
	if err != nil {
		return AnyStruct, err
	}
	return AnyStruct, nil
}
