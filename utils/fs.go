package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/utf8string"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"ravro_dcrpt/entity"
	"regexp"
	"runtime"
	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ProccesFile struct {
	NewPathFile string
	OldNamePath string
	NewNamePath string
	FileExt     string
	OldName     string
	Filename    string
	BasePath    string
	Name        string
}

func RandSeq(n int) string {
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
	err := os.Mkdir(dirName, 0775)

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

func Projectpath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return path, err
	}
	return path, nil
}

func FileProccessing(name string) (ProccesFile, error) {

	var processFile ProccesFile

	processFile.Name = name
	processFile.Filename = filepath.Base(processFile.Name)
	processFile.BasePath = filepath.Dir(processFile.Name)
	processFile.Filename = strings.Replace(processFile.Filename, ".ravro", "", 1)
	processFile.FileExt = filepath.Ext(processFile.Filename)
	processFile.OldName = fileNameWithoutExtension(processFile.Filename)
	asciiCheck := utf8string.NewString(processFile.Filename).IsASCII()
	if !asciiCheck {
		newNameRand := RandSeq(10)
		if runtime.GOOS == "windows" {
			processFile.NewPathFile = processFile.BasePath + "\\" + newNameRand + processFile.FileExt + ".ravro"
		} else {
			processFile.NewPathFile = processFile.BasePath + "/" + newNameRand + processFile.FileExt + ".ravro"
		}
		err := os.Rename(processFile.Name, processFile.NewPathFile)
		if err != nil {
			return processFile, err
		}
		processFile.Filename = newNameRand + processFile.FileExt
		processFile.Name = processFile.NewPathFile
	}
	return processFile, nil
}

func CheckDir(dirName string) bool {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return true
	}
	return false
}

func CheckPlatform(outFixpath string, process ProccesFile) ProccesFile {
	process.OldNamePath = filepath.Join(outFixpath, process.OldName+process.FileExt)
	process.NewNamePath = filepath.Join(outFixpath, process.Filename)
	return process
}

func JsonParser(process ProccesFile, AnyStruct interface{}) (interface{}, error) {

	jsonFile, err := os.Open(process.OldNamePath)
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

func ChangeDirName(reportId string, dirName string) error {

	var oldPath string
	var newPath string

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			dirNewName := filepath.Join(dirName, reportId)
			AddDir(dirNewName)
			oldPath = filepath.Join(dirName, file.Name())
			newPath = filepath.Join(dirName, reportId, file.Name())
			err := os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func ReportFiles(path, exten string) ([]string, error) {
	out, err := WalkMatch(path, exten)
	return out, err
}

func GetReportID(valuePath string) string {
	pattern := regexp.MustCompile("r[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]")
	firstMatchIndex := pattern.FindStringIndex(valuePath)
	return getSubstring(valuePath, firstMatchIndex)
}

func getSubstring(s string, indices []int) string {
	return s[indices[0]:indices[1]]
}

func CheckIsEmpty(pdf entity.Pdf) entity.Pdf {
	publicMessage := "شرح داده نشد است"
	if pdf.Report.Reproduce == "" {
		pdf.Report.Reproduce = publicMessage
	}
	if pdf.Judge.Description == "" {
		pdf.Judge.Description = publicMessage
	}
	if pdf.Judge.Vulnerability.Writeup == "" {
		pdf.Judge.Vulnerability.Writeup = publicMessage
	}
	if pdf.Judge.Vulnerability.Fix == "" {
		pdf.Judge.Vulnerability.Fix = publicMessage
	}
	if pdf.Judge.Vulnerability.Define == "" {
		pdf.Judge.Vulnerability.Define = publicMessage
	}
	if pdf.Judge.Vulnerability.Name == "" {
		pdf.Judge.Vulnerability.Name = publicMessage
	}

	return pdf
}
