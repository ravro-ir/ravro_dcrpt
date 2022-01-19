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
	report Report
	judge  Judgment
}

func main() {
	var (
		templatePath string
		outputPath   string
	)
	r := NewRequestPdf("")
	pt := ptime.Now()
	AddDir("pdfs")
	AddDir("decrypt")
	if runtime.GOOS == "windows" {
		templatePath = "template\\sample.html"
		outputPath = "pdfs\\reports.pdf"
	} else {
		templatePath = "template/sample.html"
		outputPath = "pdfs/reports.pdf"
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
		Attachment  int
	}{
		Title:     pdf.report.Title,
		PoC:       pdf.report.Description,
		CVSS:      pdf.report.CVSS,
		Reproduce: pdf.report.Reproduce,
		DateFrom:  pdf.report.DateFrom,
		Hunter:    pdf.report.HunterUsername,
		ReportID:  pdf.report.Slug,
	}
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		_, _ = r.GeneratePDF(outputPath)
		fmt.Println("[++++] pdf generated successfully")
	} else {
		fmt.Println(err)
	}
}
