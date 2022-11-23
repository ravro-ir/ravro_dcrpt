package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
	"regexp"
	"runtime"
	"strings"
)

func ReportFiles(path, exten string) ([]string, error) {
	out, err := utils.WalkMatch(path, exten)
	return out, err
}

func getSubstring(s string, indices []int) string {
	return string(s[indices[0]:indices[1]])
}

func GetReportID(valuePath string) string {
	pattern := regexp.MustCompile("r[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]")
	//welcomeMessage := "Hello guys, welcome to new york city"
	firstMatchIndex := pattern.FindStringIndex(valuePath)
	return getSubstring(valuePath, firstMatchIndex)
}

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
		lstReport, err = ReportFiles(path, "*.ravro")
		if err != nil {
			return report, err
		}
		if len(lstReport) == 0 {
			zipFile, err := ReportFiles(path, "*.zip")
			if err != nil {
				return report, err
			}
			if len(zipFile) == 0 {
				return report, err
			} else {
				CurrPath, _ := os.Getwd()
				if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
					for _, value := range zipFile {
						if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
							err := utils.Unzip(value, CurrPath+"/encrypt/"+GetReportID(value))
							if err != nil {
								fmt.Println("[----] Error : Unable to extract zip file.")
								return report, errors.New("[----] Error : Unable to extract zip file.")
							}
						} else {
							err := utils.Unzip(value, CurrPath+"\\encrypt\\"+GetReportID(value))
							if err != nil {
								fmt.Println("[----] Error : Unable to extract zip file.")
								return report, errors.New("[----] Error : Unable to extract zip file.")
							}
						}
					}

				}
			}
			lstReport, err = ReportFiles(path, "*.ravro")
			if err != nil {
				return report, err
			}
		}

		lstInfo, _ := utils.WalkMatch(path, "report_info.json")
		if len(lstReport) >= 1 {
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
	} else {
		lstReport, err = utils.WalkMatch(currentPath, "*.ravro")
		if err != nil {
			return report, err
		}
		lstReportLen := len(lstReport)
		if lstReportLen == 0 {
			return report, err
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
				_, err = utils.SslDecrypt(Process.Name, outFixpath+"/"+Process.Filename, keyFixPath)
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
