package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

type InfoReport []struct {
  InfoDescription string `json:"infoDescription"`
  InfoTitle       string `json:"infoTitle"`
  InfoSolution    string `json:"infoSolution"`
  MoreInfo        string `json:"infoMore"`
}

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
	Attachment      bool
	Scenario        string `json:"scenario"`
  ReportInfo      InfoReport
}

func DcrptReport(currentPath, keyFixPath, outFixpath string, checkStatus bool) (Report, error) {
	var report Report
  var infoReport InfoReport
 	var (
		path      string
		err       error
		lstReport []string
	)
  
	if currentPath == "" {
		path, err = projectpath()
		if err != nil {
			return report, err
		}
		lstReport, err = WalkMatch(path, "*.ravro")
		if err != nil {
			return report, err
		}
		lstReportLen := len(lstReport)
		if lstReportLen == 0 {
			return report, err
		}
		if lstReportLen > 1 {
			report.Attachment = true
		}
    lstInfo, _ := WalkMatch(path, "*.json")
    
    jsonFile, err := os.Open(lstInfo[0])
	  reportValue, _ := ioutil.ReadAll(jsonFile)
 	  err = json.Unmarshal(reportValue, &infoReport)
    if err != nil {
		  return report, err
	  }
	  err = jsonFile.Close()
	  if err != nil {
		  return report, err
    }
    report.ReportInfo = infoReport
	} else {
		lstReport, err = WalkMatch(currentPath, "*.ravro")
		if err != nil {
			return report, err
		}
		lstReportLen := len(lstReport)
		if lstReportLen == 0 {
			return report, err
		}
		if lstReportLen > 1 {
			report.Attachment = true
		}

	}
	for _, name := range lstReport {
		if runtime.GOOS == "windows" {
			if !checkStatus {
				if !strings.Contains(name, "\\encrypt\\") {
					continue
				}
			}
			if !strings.Contains(name, "\\report\\") {
				continue
			}
		} else {
			if !strings.Contains(name, "encrypt/") {
				continue
			}
			if !strings.Contains(name, "/report/") {
				continue
			}
		}
		Process, err := fileProccessing(name)
		if err != nil {
			return report, err
		}
		if runtime.GOOS == "windows" {
			_, err = SslDecrypt(Process.name, outFixpath+"\\"+Process.filename, keyFixPath)
		} else {
			if strings.Contains(outFixpath, "/") {
				_, err = SslDecrypt(Process.name, outFixpath+Process.filename, keyFixPath)
			} else {
				_, err = SslDecrypt(Process.name, outFixpath+"/"+Process.filename, keyFixPath)
			}
		}
		if err != nil {
			return report, err
		}
		process := CheckPlatform(outFixpath, Process)
		err = os.Rename(process.newNamePath, process.oldNamePath)
		if err != nil {
			return report, err
		}
		if strings.Index(process.oldName, "data") != 0 {
			continue
		}
		_, err = JsonParser(process, &report)
		if err = os.Remove(process.oldNamePath); err != nil {
			log.Fatal(err)
		}
    //jsonFile, err := os.Open(currentPath + "report_info.json")
	  //reportValue, _ := ioutil.ReadAll(jsonFile)
	  //err = json.Unmarshal(reportValue, &report)
    //if err != nil {
		// return report, err
	  //}
	  //err = jsonFile.Close()
	  //if err != nil {
		//return report, err
	  //}
	}
	return report, nil
}
