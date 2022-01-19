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

type Amendment struct {
	AttachmentsCount int    `json:"attachmentsCount"`
	Description      string `json:"description"`
	HunterUsername   string `json:"hunterUsername"`
	CompanyUsername  string `json:"companyUsername"`
	SubmissionDate   string `json:"submissionDate"`
}

func DcrptAmendment() ([]string, error) {
	var amendment Amendment
	var lstMore []string
	path, err := os.Getwd()
	if err != nil {
		return lstMore, err
	}
	lstAmendment, _ := WalkMatch(path, "*.ravro")
	for _, name := range lstAmendment {
		var NewPathFile string
		var oldNamePath string
		var newNamePath string
		if !strings.Contains(name, "amendment-") {
			continue
		}
		newName := strings.ReplaceAll(name, " ", "")
		err := os.Rename(name, newName)
		if err != nil {
			return lstMore, err
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
				return lstMore, err
			}
			filename = newNameRand + fileExt
			name = NewPathFile
		}
		_, err = SslDecrypt(name, filename)
		if err != nil {
			return lstMore, err
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
			return lstMore, err
		}
		if !strings.Contains(oldNamePath, "data") {
			continue
		}

		jsonFile, err := os.Open(oldNamePath)
		reportValue, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(reportValue, &amendment)
		if err != nil {
			return lstMore, err
		}
		err = jsonFile.Close()
		if err != nil {
			return lstMore, err
		}
		if err = os.Remove(oldNamePath); err != nil {
			log.Fatal(err)
		}
		lstMore = append(lstMore, amendment.Description)
		//return lstMore, nil
	}
	return lstMore, nil
}
