package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gomarkdown/markdown"
	ptime "github.com/yaa110/go-persian-calendar"
	"io/ioutil"
	"log"
	"os"
	"ravro_dcrpt/core"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	rvrVersion    = "v1.0.1"
	publicMessage = "شرح داده نشده است."
	noMsg         = "ثبت نشد"
)

func main() {

	var (
		templatePath string
		outputPath   string
		keyFixPath   string
		outFixpath   string
		curretnPath  string
		status       bool
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
	inputDir := flag.String("in", "in", "input directory of report encrypt file")
	version := flag.String("version", ">> Current Version : ravro_dcrpt/1.0.2", "")
	homePage := flag.String("homepage", ">> Github : https://github.com/ravro-ir/ravro_dcrp", "")
	issue := flag.String("issue", ">> Issue : https://github.com/ravro-ir/ravro_dcrp/issues", "")
	author := flag.String("author", ">> Author : Ravro Development Team (RDT)", "")
	help := flag.String("help", ">> Help : ravro_dcrpt --help \n\n", "")
	latest := flag.Bool("latest", true, "")

	outputDir := flag.String("out", "out", "output directory for decrypt report file ")
	key := flag.String("key", "key", "key.private")
	init := flag.String("init", "", "input directory of report encrypt file")
	update := flag.Bool("udate", false, "Update ravro decryptor")
	format := flag.Bool("json", false, "Convert report to json")
	log := flag.Bool("log", false, "Store logs in log.txt")

	flag.Parse()
	fmt.Println(*version)
	fmt.Println(*homePage)
	fmt.Println(*issue)
	fmt.Println(*author)
	fmt.Println(*help)

	if *latest {
		ver, err := utils.LatestVersion()
		if err != nil {
			LogCheck(*log, err)
		} else {
			newVer := fmt.Sprintf("ravro_dcrpt/%s", ver)
			if rvrVersion != ver {
				fmt.Println(fmt.Sprintf("\n New version (%s) released, "+
					"Please use command : ./ravro_dcrpt -u \n", newVer))
				return
			}
		}
	}
	if *update {
		fmt.Println("[++++] Downloading latest version")
		err := utils.HttpGet()
		if err != nil {
			LogCheck(*log, err)
			fmt.Println("[----] Unable to download latest file, Maybe your internet connection is bad," +
				"or please usage : `./ravro_dcrpt -log` for see monitoring error message")
		}
		// Extract zip file
		return
	}
	if *init == "init" {
		utils.AddDir("decrypt")
		utils.AddDir("encrypt")
		utils.AddDir("key")
		fmt.Println("[++] Created directory decrypt, encrypt, key")
		return
	} else {
		for i := range lstDir {
			if *inputDir != "in" && *outputDir != "out" && *key != "key" {
				break
			}
			out := utils.CheckDir(lstDir[i])
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
		outFixpath = *outputDir
	}
	if *key != "key" {
		status = true
		keyFixPath = *key
	}
	r := utils.NewRequestPdf("")
	pt := ptime.Now()

	fmt.Println("[++++] Starting for decrypting Report . . . ")
	report, err := core.DcrptReport(curretnPath, keyFixPath, outFixpath, status)
	if err != nil {
		LogCheck(*log, err)
		fmt.Println("[----] Error : The private key is incorrect")
		return
	}
	if report.Title == "" {
		LogCheck(*log, err)
		fmt.Println("[----] The input file for decryption is not correct.")
		return
	}
	fmt.Println("[++++] Starting for decrypting Judgment . . . ")
	judge, err := core.DcrptJudgment(curretnPath, keyFixPath, outFixpath, status)
	if err != nil {
		LogCheck(*log, err)
		fmt.Println(err)
		return
	}
	fmt.Println("[++++] Starting for decrypting Amendment . . . ")
	amendment, err := core.DcrptAmendment(curretnPath, keyFixPath, outFixpath)
	if err != nil {
		LogCheck(*log, err)
		fmt.Println(err)
		return
	}
	utils.AddDir("template")
	utils.HtmlTemplate(templatePath)
	moreInfo := strings.Join(amendment[:], ",")
	if moreInfo == "" {
		moreInfo = publicMessage
	}
	if report.Reproduce == "" {
		report.Reproduce = publicMessage
	}
	dateTo := strconv.Itoa(pt.Year()) + "/" + strconv.Itoa(int(pt.Month())) + "/" + strconv.Itoa(pt.Day())
	pdf := entity.Pdf{Judge: judge, Report: report}

	dateFrom, outputPath := Validate(report, outputPath, pdf)
	fmt.Println("[++++] Starting report to pdf . . . ")
	if *format {
		file, _ := json.MarshalIndent(pdf.Judge, "", " ")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			_ = ioutil.WriteFile("decrypt//juror.json", file, 0644)
		} else {
			_ = ioutil.WriteFile("decrypt\\juror.json", file, 0644)
		}
		reportd, _ := json.MarshalIndent(pdf.Report, "", " ")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			_ = ioutil.WriteFile("decrypt//repo.json", reportd, 0644)
		} else {
			_ = ioutil.WriteFile("decrypt\\repo.json", reportd, 0644)
		}
		amendments, _ := json.MarshalIndent(pdf.Amendment, "", " ")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			_ = ioutil.WriteFile("decrypt//moreinfo.json", amendments, 0644)
		} else {
			_ = ioutil.WriteFile("decrypt\\moreinfo.json", amendments, 0644)
		}
	}

	md := []byte(pdf.Report.Description)
	templateData := TemplateStruct(md, pdf, dateFrom, dateTo, moreInfo, report)
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		_, err = r.GeneratePDF(outputPath)
		if err != nil {
			LogCheck(*log, err)
			fmt.Println("[----] failed to remove html template,")
		}
		err := os.RemoveAll("template")
		if err != nil {
			LogCheck(*log, err)
			fmt.Println("[----] failed to remove html template,")
		}
		fmt.Println("[++++] PDF generated successfully")
		err = utils.ChangeDirName(report.Slug, outFixpath)
		if err != nil {
			LogCheck(*log, err)
		}
	} else {
		LogCheck(*log, err)
		fmt.Println(err)
	}
}

func AttachmentFiles(info entity.InfoReport) string {
	var attach string
	if len(info.Details.Attachments) == 0 {
		return "<tr>\n<td>فایلی توسط شکارچی پیوست نشده است</td>\n</tr>"
	}
	for _, content := range info.Details.Attachments {
		attach += fmt.Sprintf("<tr>\n<td>%s</td>\n</tr>", content.Filename)
	}
	return attach
}

func JugmentUser(info entity.InfoReport) string {
	var judge string
	for _, content := range info.Details.Judges {
		judge += content.Name + "  "
	}
	return judge
}

func ConString(info entity.InfoReport) string {
	var infoMore string
	for _, content := range info.Tags {
		infoMore += " عنوان گزارش : " + content.Title + "<br />"
		infoMore += " توضیحات گزارش : " + "<br />" + content.Description + "<br />"
		infoMore += "راه حل : " + "<br />" + content.Solution + "<br />"
		infoMore += "اطلاعات بیشتر : " + "<br />" + content.MoreInfo + "<br />"
		infoMore += "<hr>"
	}
	return infoMore
}

func Validate(report entity.Report, outputPath string, pdf entity.Pdf) (string, string) {
	var dateSubmit string
	if report.CompanyUsername == "" {
		outputPath = strings.Replace(outputPath, "reports", utils.RandSeq(8), 1)
	} else {
		outputPath = strings.Replace(outputPath, "reports", report.CompanyUsername+"__"+report.Slug+"__"+report.HunterUsername, 1)
	}
	pdf = utils.CheckIsEmpty(pdf)

	fmt.Println("[++++] Decrypted successfully ")
	if pdf.Report.SubmissionDate == "" {
		dateSubmit = pdf.Report.DateFrom
	} else {
		dateSubmit = pdf.Report.SubmissionDate
	}
	dateSubmited := strings.Split(dateSubmit, " ")
	if len(dateSubmited) == 0 {
		log.Fatalln("[----] Error : The submit date is empty, we think your report path is incorrect, (Valid Path: (encrypt/ir2022-01-10-0001))")
	}
	dateReport := strings.Split(dateSubmited[0], "-")
	dateFrom := MiladiToShamsi(ArrayStringToInt(dateReport[0]), ArrayStringToInt(dateReport[1]), ArrayStringToInt(dateReport[2]))
	return dateFrom, outputPath
}

func ArrayStringToInt(date string) int {
	newData, err := strconv.Atoi(date)
	if err != nil {
		log.Fatal(err)
	}
	return newData
}

func MiladiToShamsi(year, month, day int) string {
	pt := ptime.Now()
	var t = time.Date(year, time.Month(month), day, 12, 1, 1, 0, ptime.Iran())
	pt = ptime.New(t)
	shamsiDate := pt.Format("yyyy/MM/dd")
	return shamsiDate
}

func TemplateStruct(md []byte, pdf entity.Pdf, dateFrom, dateTo, moreInfo string, report entity.Report) any {
	output := markdown.ToHTML(md, nil, nil)
	infoTrain := ConString(report.ReportInfo)
	dateOf := strings.Split(report.DateFrom, "-")
	if len(dateOf) == 0 {
		log.Fatalln("[----] Error : The datefrom is invalid,")
	}
	report.DateFrom = MiladiToShamsi(ArrayStringToInt(dateOf[0]), ArrayStringToInt(dateOf[1]), ArrayStringToInt(dateOf[2]))

	dateTwo := strings.Split(report.DateTo, "-")
	if len(dateTwo) == 0 {
		log.Fatalln("[----] Error : The datefrom is invalid,")
	}
	report.DateTo = MiladiToShamsi(ArrayStringToInt(dateTwo[0]), ArrayStringToInt(dateTwo[1]), ArrayStringToInt(dateTwo[2]))
	if infoTrain == "" {
		infoTrain = publicMessage
	}
	if report.ReportInfo.Details.Cvss.Judge.Score == "" {
		report.ReportInfo.Details.Cvss.Judge.Score = noMsg
	}
	if report.ReportInfo.Details.Cvss.Hunter.Score == "" {
		report.ReportInfo.Details.Cvss.Hunter.Score = noMsg
	}
	if report.ReportInfo.Details.Cvss.Hunter.Vector == "" {
		report.ReportInfo.Details.Cvss.Hunter.Vector = noMsg
	}
	if report.ReportInfo.Details.Cvss.Judge.Vector == "" {
		report.ReportInfo.Details.Cvss.Judge.Vector = noMsg
	}
	templateData := struct {
		Title           string
		Description     string
		PoC             string
		DateFrom        string
		CVSSJudge       string
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
		JudgeUser       string
		ScoreJudge      string
		ScoreHunter     string
		CVSSHunter      string
		RangeDate       string
		Targets         string
		Status          string
	}{
		Title:           pdf.Report.Title,
		PoC:             string(output),
		CVSSJudge:       report.ReportInfo.Details.Cvss.Judge.Vector,
		Reproduce:       pdf.Report.Reproduce,
		DateFrom:        dateFrom,
		Hunter:          pdf.Report.HunterUsername,
		ReportID:        pdf.Report.Slug,
		Amount:          pdf.Judge.Reward,
		JudgeInfo:       pdf.Judge.Description,
		VulDefine:       pdf.Judge.Vulnerability.Define,
		VulType:         pdf.Judge.Vulnerability.Name,
		VulWriteup:      pdf.Judge.Vulnerability.Writeup,
		VulFix:          pdf.Judge.Vulnerability.Fix,
		DateTo:          dateTo,
		MoreInfo:        moreInfo,
		CompanyUserName: pdf.Report.CompanyUsername,
		Ips:             pdf.Report.Ips,
		Attachment:      AttachmentFiles(report.ReportInfo),
		Scenario:        pdf.Report.Scenario,
		LinkMoreInfo:    infoTrain,
		RavroVer:        rvrVersion,
		JudgeUser:       JugmentUser(report.ReportInfo),
		ScoreJudge:      report.ReportInfo.Details.Cvss.Judge.Score,
		ScoreHunter:     report.ReportInfo.Details.Cvss.Hunter.Score,
		CVSSHunter:      report.ReportInfo.Details.Cvss.Hunter.Vector,
		RangeDate:       "از" + report.DateFrom + "  تا " + report.DateTo,
		Targets:         report.ReportInfo.Details.Target,
		Status:          report.ReportInfo.Details.CurrentStatus,
	}
	return templateData
}

func LogCheck(status bool, Msg any) {
	if status {
		core.ErrorLogger.Println(Msg)
	}
}
