package main

import (
	"fmt"
	ptime "github.com/yaa110/go-persian-calendar"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Pdf struct {
	report    Report
	judge     Judgment
	amendment Amendment
}

func main() {
	var (
		templatePath string
		outputPath   string
	)
	r := NewRequestPdf("")
	pt := ptime.Now()
	AddDir("decrypt")
	if runtime.GOOS == "windows" {
		templatePath = "template\\sample.html"
		outputPath = "decrypt\\reports.pdf"
	} else {
		templatePath = "template/sample.html"
		outputPath = "decrypt/reports.pdf"
	}
	fmt.Println("[++++] Starting for decrypting . . . ")
	judge, err := DcrptJudgment()
	if err != nil {
		log.Fatalln("Error judge")
	}
	report, err := DcrptReport()
	if err != nil {
		log.Fatalln("Error report")
	}
	amendment, err := DcrptAmendment()
	if err != nil {
		log.Fatalln("Error amendment")
	}
	moreinfo := strings.Join(amendment[:], ",")
	//if (Amendment{}) != amendment  {
	//	fmt.Println("is zero value")
	//}
	dateTo := strconv.Itoa(pt.Year()) + "/" + strconv.Itoa(int(pt.Month())) + "/" + strconv.Itoa(pt.Day())
	pdf := Pdf{judge: judge, report: report}
	fmt.Println("[++++] decrypted successfully ")
	fmt.Println("[++++] Starting report to pdf . . . ")
	dataFrom := pdf.report.DateFrom
	date := strings.Split(dataFrom, "-")
	year, _ := strconv.Atoi(date[0])
	month, _ := strconv.Atoi(date[1])
	day, _ := strconv.Atoi(date[2])
	var t time.Time = time.Date(year, time.Month(month), day, 12, 1, 1, 0, ptime.Iran())
	pt = ptime.New(t)
	dataFrom = pt.Format("yyyy/MM/dd")
	templateData := struct {
		Title       string
		Description string
		PoC         string
		DateFrom    string
		CVSS        string
		Reproduce   string
		Hunter      string
		ReportID    string
		Amount      int
		Score       int
		JudgeInfo   string
		DateTo      string
		MoreInfo    string
	}{
		Title:     pdf.report.Title,
		PoC:       pdf.report.Description,
		CVSS:      pdf.judge.Cvss.Value,
		Reproduce: pdf.report.Reproduce,
		DateFrom:  pdf.report.DateFrom,
		Hunter:    pdf.report.HunterUsername,
		ReportID:  pdf.report.Slug,
		Amount:    pdf.judge.Reward,
		JudgeInfo: pdf.judge.Description,
		DateTo:    dateTo,
		MoreInfo:  moreinfo,
	}
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		_, _ = r.GeneratePDF(outputPath)
		fmt.Println("[++++] pdf generated successfully")
	} else {
		fmt.Println(err)
	}
}
