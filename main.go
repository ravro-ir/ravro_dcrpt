package main

import (
	"encoding/json"
	"fmt"
	ptime "github.com/yaa110/go-persian-calendar"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
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
	Attachment      int    `json:"attachmentsCount"`
}

func AddDir(name string) {
	if err := ensureDir(name); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}
}

func main() {
	var (
		report       Report
		templatePath string
		outputPath   string
		jsonVulPath  string
	)
	r := NewRequestPdf("")
	pt := ptime.Now()
	AddDir("pdfs")
	AddDir("datadecrypt")
	if runtime.GOOS == "windows" {
		templatePath = "template\\sample.html"
		outputPath = "pdfs\\reports.pdf"
		jsonVulPath = "datadecrypt\\data"
	} else {
		templatePath = "template/sample.html"
		outputPath = "pdfs/reports.pdf"
		jsonVulPath = "datadecrypt/data"
	}
	fmt.Println("[++++] Starting for decrypting . . . ")
	ReadCurrentDir()
	fmt.Println("[++++] decrypted successfully ")
	fmt.Println("[++++] Starting report to pdf . . . ")
	jsonFile, err := os.Open(jsonVulPath)
	byteValueReport, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValueReport, &report)
	if err != nil {
		return
	}
	dataFrom := report.DateFrom
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
		Title:      report.Title,
		PoC:        report.Description,
		CVSS:       report.CVSS,
		Reproduce:  report.Reproduce,
		DateFrom:   report.DateFrom,
		Hunter:     report.HunterUsername,
		ReportID:   report.Slug,
		Attachment: report.Attachment,
	}
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		_, _ = r.GeneratePDF(outputPath)
		fmt.Println("[++++] pdf generated successfully")
	} else {
		fmt.Println(err)
	}
}
