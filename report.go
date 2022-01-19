package main

import (
	"encoding/json"
	"golang.org/x/exp/utf8string"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Report struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Reproduce       string `json:"reproduce"`
	DateFrom        string `json:"dateFrom"`
	CVSS            string `json:"cvss"`
	HunterUsername  string `json:"hunterUsername"`
	CompanyUsername string `json:"companyUsername"`
	Slug            string `json:"slug"`
}

func DcrptReport() (Report, error) {
	var report Report
	path, err := os.Getwd()
	if err != nil {
		return report, err
	}
	lstReport, _ := WalkMatch(path, "*.ravro")
	for _, name := range lstReport {
		var NewPathFile string
		var oldNamePath string
		var newNamePath string
		if !strings.Contains(name, "\\report\\") {
			continue
		}
		newName := strings.ReplaceAll(name, " ", "")
		err := os.Rename(name, newName)
		if err != nil {
			return report, err
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
				return report, err
			}
			filename = newNameRand + fileExt
			name = NewPathFile
		}
		_, err = SslDecrypt(name, filename)
		if err != nil {
			return report, err
		}
		if runtime.GOOS == "windows" {
			oldNamePath = "decrypt" + "\\" + oldName + fileExt
			newNamePath = "decrypt" + "\\" + filename
		} else {
			oldNamePath = "decrypt" + "/" + oldName + fileExt
			newNamePath = "decrypt" + "/" + filename
		}
		err = os.Rename(newNamePath, oldNamePath)
		if err != nil {
			return report, err
		}
		if !strings.Contains(oldNamePath, "data") {
			continue
		}
		jsonFile, err := os.Open(oldNamePath)
		byteValueReport, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(byteValueReport, &report)
		if err != nil {
			return report, err
		}
		jsonFile.Close()
		if err = os.Remove(oldNamePath); err != nil {
			log.Fatal(err)
		}
		return report, nil
	}
	return report, nil
}
