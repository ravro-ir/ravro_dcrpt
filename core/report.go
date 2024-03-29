package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
	"runtime"
	"strings"
)

const jsonInfo string = "report_info.json"

func DcrptReport(currentPath, keyFixPath, outFixpath string, checkStatus bool) (entity.Report, error) {
	var report entity.Report
	var infoReport entity.InfoReport
	var (
		err       error
		lstReport []string
	)
	lstInfo, _ := utils.WalkMatch(currentPath, jsonInfo)
	if len(lstInfo) >= 1 {
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
	}
	lstReport, err = utils.WalkMatch(currentPath, "*.ravro")
	if err != nil {
		return report, err
	}
	lstReportLen := len(lstReport)
	if lstReportLen == 0 {
		return report, err
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
			if !strings.Contains(name, "/report/") {
				continue
			}
		}
		Process, err := utils.FileProccessing(name)
		if err != nil {
			return report, err
		}
		reportPath := filepath.Join(outFixpath, Process.Filename)
		_, err = utils.SslDecrypt(Process.Name, reportPath, keyFixPath)
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
