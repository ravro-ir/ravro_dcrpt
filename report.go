package main

import (
	"log"
	"os"
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
	SubmissionDate  string `json:"submissionDate"`
	Ips             string `json:"ips"`
}

func DcrptReport() (Report, error) {
	var report Report
	path, err := projectpath()
	if err != nil {
		return report, err
	}
	lstReport, _ := WalkMatch(path, "*.ravro")
	for _, name := range lstReport {
		if runtime.GOOS == "windows" {
			if !strings.Contains(name, "\\report\\") {
				continue
			}
		} else {
			if !strings.Contains(name, "/report/") {
				continue
			}
		}
		Process, err := fileProccessing(name)
		if err != nil {
			return report, err
		}
		_, err = SslDecrypt(Process.name, Process.filename)
		if err != nil {
			return report, err
		}
		process := CheckPlatform(Process)
		err = os.Rename(process.newNamePath, process.oldNamePath)
		if err != nil {
			return report, err
		}
		if !strings.Contains(process.oldNamePath, "data") {
			continue
		}
		_, err = JsonParser(process, &report)
		if err = os.Remove(process.oldNamePath); err != nil {
			log.Fatal(err)
		}
	}
	return report, nil
}
