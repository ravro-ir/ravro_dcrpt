package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
	"runtime"
	"strings"
)

func DcrptReport(currentPath, keyFixPath, outFixpath string, checkStatus bool) (entity.Report, error) {
	var report entity.Report
	var infoReport entity.InfoReport
	var (
		path      string
		err       error
		lstReport []string
	)

	if currentPath == "" {
		path, err = utils.Projectpath()
		if err != nil {
			return report, err
		}
		lstReport, err = utils.WalkMatch(path, "*.ravro")
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
		lstInfo, _ := utils.WalkMatch(path, "report_info.json")
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
		lstReport, err = utils.WalkMatch(currentPath, "*.ravro")
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
		Process, err := utils.FileProccessing(name)
		if err != nil {
			return report, err
		}
		if runtime.GOOS == "windows" {
			_, err = utils.SslDecrypt(Process.Name, outFixpath+"\\"+Process.Filename, keyFixPath)
		} else {
			if strings.Contains(outFixpath, "/") {
				_, err = utils.SslDecrypt(Process.Name, outFixpath+Process.Filename, keyFixPath)
			} else {
				_, err = utils.SslDecrypt(Process.Name, outFixpath+"/"+Process.Filename, keyFixPath)
			}
		}
		if err != nil {
			return report, err
		}
		process := utils.CheckPlatform(outFixpath, Process)
		err = os.Rename(process.NewNamePath, process.OldNamePath)
		if err != nil {
			return report, err
		}
		if strings.Index(process.OldName, "data") != 0 {
			continue
		}
		_, err = utils.JsonParser(process, &report)
		if err = os.Remove(process.OldNamePath); err != nil {
			log.Fatal(err)
		}
	}
	return report, nil
}
