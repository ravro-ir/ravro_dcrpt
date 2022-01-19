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

type JudgeCvss struct {
	Value  string `json:"value"`
	Rating string `json:"rating"`
}

type Judgment struct {
	Score       int    `json:"score"`
	Reward      int    `json:"reward"`
	Description string `json:"description"`
	Cvss        JudgeCvss
}

func DcrptJudgment() (Judgment, error) {
	var judgment Judgment
	path, err := os.Getwd()
	if err != nil {
		return judgment, err
	}
	lstJudge, _ := WalkMatch(path, "*.ravro")
	for _, name := range lstJudge {
		var NewPathFile string
		var oldNamePath string
		var newNamePath string
		if !strings.Contains(name, "\\judgment\\") {
			continue
		}
		newName := strings.ReplaceAll(name, " ", "")
		err := os.Rename(name, newName)
		if err != nil {
			return judgment, err
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
				return judgment, err
			}
			filename = newNameRand + fileExt
			name = NewPathFile
		}
		_, err = SslDecrypt(name, filename)
		if err != nil {
			return judgment, err
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
			return judgment, err
		}
		jsonFile, err := os.Open(oldNamePath)
		judgmentValue, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(judgmentValue, &judgment)
		if err != nil {
			return judgment, err
		}
		err = jsonFile.Close()
		if err != nil {
			return Judgment{}, err
		}
		if err = os.Remove(oldNamePath); err != nil {
			log.Fatal(err)
		}
		return judgment, nil
	}
	return judgment, nil

}
