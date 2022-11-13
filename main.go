package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	ptime "github.com/yaa110/go-persian-calendar"
)

type Pdf struct {
	report    Report
	judge     Judgment
	amendment Amendment
}

const rvrVersion = "v1.0.2"

func main() {

	var (
		templatePath string
		outputPath   string
		keyFixPath   string
		outFixpath   string
		curretnPath  string
		status       bool
		AttachStatus string
	)
	lstDir := []string{"encrypt", "decrypt", "key"}
	if runtime.GOOS == "windows" {
		templatePath = "template\\sample.html"
		outputPath = "decrypt\\reports.pdf"
		keyFixPath = "key/key.private"
		outFixpath = "decrypt"
	} else {
		templatePath = "template/sample.html"
		outputPath = "decrypt/reports.pdf"
		keyFixPath = "key/key.private"
		outFixpath = "decrypt"
	}
	// تاریخ بررسی - اعضای تیم داوری - تگ
	publicMessage := "شرح داده نشده است."
	inputDir := flag.String("in", "in", "input directory of report encrypt file")
	version := flag.String("version", ">> Version : ravro_dcrpt/1.0.1", "")
	homePage := flag.String("homepage", ">> HomePage: https://github.com/ravro-ir/ravro_dcrp", "")
	issue := flag.String("issue", ">> Issue: https://github.com/ravro-ir/ravro_dcrp/issues", "")
	author := flag.String("author", ">> Author : Ramin Farajpour Cami", "")
	help := flag.String("help", ">> Help : ravro_dcrpt --help \n\n", "")

	outputDir := flag.String("out", "out", "output directory for decrypt report file ")
	key := flag.String("key", "key", "key.private")
	init := flag.String("init", "", "input directory of report encrypt file")
	update := flag.String("update", "", "Update ravro decryptor: update=yes")

	flag.Parse()
	fmt.Println(*version)
	fmt.Println(*homePage)
	fmt.Println(*issue)
	fmt.Println(*author)
	fmt.Println(*help)

	if *update == "yes" {
		fmt.Println("[++++] Downloading latest version")
		HttpGet()
		return
	}
	if *init == "init" {
		AddDir("decrypt")
		AddDir("encrypt")
		AddDir("key")
		fmt.Println("[++] Created directory decrypt, encrypt, key")
		return
	} else {
		for i := range lstDir {
			if *inputDir != "in" && *outputDir != "out" && *key != "key" {
				break
			}
			out := CheckDir(lstDir[i])
			if out {
				fmt.Println("[---] encrypt, decrypt, key is not exist.")
				fmt.Println("[+++] Usage : ravro_dcrpt -init=init")
				return
			}
		}
	}
	status = false
	if *inputDir != "in" {
		status = true
		curretnPath = *inputDir
	}
	if *outputDir != "out" {
		status = true
		outputPath = *outputDir + "reports.pdf"
		//path, err := os.Getwd()
		//if err != nil {
		//	log.Println(err)
		//}
		outFixpath = *outputDir
	}
	if *key != "key" {
		status = true
		keyFixPath = *key
	}
	r := NewRequestPdf("")
	pt := ptime.Now()

	fmt.Println("[++++] Starting for decrypting Report . . . ")
	report, err := DcrptReport(curretnPath, keyFixPath, outFixpath, status)
	if err != nil {
		fmt.Println(err)
		log.Fatalln("[----] Error : The private key is incorrect")
	}
	if report.Title == "" {
		fmt.Println("[----] The input file for decryption is not correct.")
		return
	}
	if report.Attachment {
		AttachStatus = "دارد"
	} else {
		AttachStatus = "ندارد"
	}
	fmt.Println("[++++] Starting for decrypting Judgment . . . ")
	judge, err := DcrptJudgment(curretnPath, keyFixPath, outFixpath, status)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("[++++] Starting for decrypting Amendment . . . ")
	amendment, err := DcrptAmendment(curretnPath, keyFixPath, outFixpath)
	if err != nil {
		log.Fatalln(err)
	}

	AddDir("template")
	HtmlTemplate(templatePath)
	moreInfo := strings.Join(amendment[:], ",")
	if moreInfo == "" {
		moreInfo = publicMessage
	}
	dateTo := strconv.Itoa(pt.Year()) + "/" + strconv.Itoa(int(pt.Month())) + "/" + strconv.Itoa(pt.Day())
	pdf := Pdf{judge: judge, report: report}

	dateFrom, outputPath := Validate(report, outputPath, pdf)
	fmt.Println("[++++] Starting report to pdf . . . ")

	md := []byte(pdf.report.Description)
	templateData := TemplateStruct(md, pdf, dateFrom, dateTo, moreInfo, AttachStatus, report)
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		_, _ = r.GeneratePDF(outputPath)
		err := os.RemoveAll("template")
		if err != nil {
			fmt.Println("[----] failed to remove html template,")
		}
		fmt.Println("[++++] pdf generated successfully")
		ChangeDirName(report.Slug, outFixpath)
	} else {
		fmt.Println(err)
	}
}

func ConString(info InfoReport) string {
	var infoMore string
	for _, content := range info {
		infoMore += " عنوان گزارش : " + content.InfoTitle + "<br />"
		infoMore += " توضیحات گزارش : " + "<br />" + content.InfoDescription + "<br />"
		infoMore += "راه حل : " + "<br />" + content.InfoSolution + "<br />"
		infoMore += "اطلاعات بیشتر : " + "<br />" + content.MoreInfo + "<br />"
		infoMore += "<hr>"
	}
	return infoMore
}

func Validate(report Report, outputPath string, pdf Pdf) (string, string) {
	var dateSubmit string
	pt := ptime.Now()
	if report.CompanyUsername == "" {
		outputPath = strings.Replace(outputPath, "reports", randSeq(8), 1)
	} else {
		outputPath = strings.Replace(outputPath, "reports", report.CompanyUsername+"__"+report.Slug+"__"+report.HunterUsername, 1)
	}
	pdf = CheckIsEmpty(pdf)

	fmt.Println("[++++] decrypted successfully ")
	if pdf.report.SubmissionDate == "" {
		dateSubmit = pdf.report.DateFrom
	} else {
		dateSubmit = pdf.report.SubmissionDate
	}
	dateSubmited := strings.Split(dateSubmit, " ")
	if len(dateSubmited) == 0 {
		log.Fatalln("[----] Error : The submit date is empty, we think your report path is incorrect, (Valid Path: (encrypt/ir2022-01-10-0001))")
	}
	dateReport := strings.Split(dateSubmited[0], "-")
	year, err := strconv.Atoi(dateReport[0])
	if err != nil {
		log.Fatalln(err)
	}
	month, err := strconv.Atoi(dateReport[1])
	if err != nil {
		log.Fatalln(err)
	}
	day, err := strconv.Atoi(dateReport[2])
	if err != nil {
		log.Fatal(err)
	}
	var t = time.Date(year, time.Month(month), day, 12, 1, 1, 0, ptime.Iran())
	pt = ptime.New(t)
	dateFrom := pt.Format("yyyy/MM/dd")

	return dateFrom, outputPath
}

func TemplateStruct(md []byte, pdf Pdf, dateFrom, dateTo, moreInfo, AttachStatus string, report Report) any {
	output := markdown.ToHTML(md, nil, nil)
	templateData := struct {
		Title           string
		Description     string
		PoC             string
		DateFrom        string
		CVSS            string
		Reproduce       string
		Hunter          string
		ReportID        string
		Amount          int
		Score           int
		JudgeInfo       string
		DateTo          string
		MoreInfo        string
		CompanyUserName string
		Ips             string
		VulDefine       string
		VulType         string
		VulFix          string
		VulWriteup      string
		Attachment      string
		Scenario        string
		LinkMoreInfo    string
		RavroVer        string
	}{
		Title:           pdf.report.Title,
		PoC:             string(output),
		CVSS:            pdf.judge.Cvss.Value,
		Reproduce:       pdf.report.Reproduce,
		DateFrom:        dateFrom,
		Hunter:          pdf.report.HunterUsername,
		ReportID:        pdf.report.Slug,
		Amount:          pdf.judge.Reward,
		JudgeInfo:       pdf.judge.Description,
		VulDefine:       pdf.judge.Vulnerability.Define,
		VulType:         pdf.judge.Vulnerability.Name,
		VulWriteup:      pdf.judge.Vulnerability.Writeup,
		VulFix:          pdf.judge.Vulnerability.Fix,
		DateTo:          dateTo,
		MoreInfo:        moreInfo,
		CompanyUserName: pdf.report.CompanyUsername,
		Ips:             pdf.report.Ips,
		Attachment:      AttachStatus,
		Scenario:        pdf.report.Scenario,
		LinkMoreInfo:    ConString(report.ReportInfo),
		RavroVer:        rvrVersion,
	}
	return templateData
}
