package main

import (
	"os"
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

func DcrptAmendment(currentPath, keyFixPath, outFixPath string) ([]string, error) {
	var amendment Amendment
	var (
		path    string
		err     error
		lstMore []string
	)
	if currentPath == "" {
		path, err = projectpath()
		if err != nil {
			return lstMore, err
		}
	}

	lstAmendment, _ := WalkMatch(path, "*.ravro")
	for _, name := range lstAmendment {

		if !strings.Contains(name, "amendment-") {
			continue
		}
		Process, err := fileProccessing(name)
		if err != nil {
			return lstMore, err
		}
		if runtime.GOOS == "windows" {
			_, err = SslDecrypt(Process.name, outFixPath+"\\"+Process.filename, keyFixPath)
		} else {
			_, err = SslDecrypt(Process.name, outFixPath+"/"+Process.filename, keyFixPath)
		}
		if err != nil {
			return lstMore, err
		}
		process := CheckPlatform(outFixPath, Process)
		err = os.Rename(process.newNamePath, process.oldNamePath)
		if err != nil {
			return lstMore, err
		}
		if !strings.Contains(process.oldNamePath, "data") {
			continue
		}
		_, err = JsonParser(process, &amendment)
		if err = os.Remove(process.oldNamePath); err != nil {
			return lstMore, err
		}
		lstMore = append(lstMore, amendment.Description)
	}
	return lstMore, nil
}
