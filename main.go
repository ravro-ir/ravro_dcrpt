package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gomarkdown/markdown"
	"github.com/manifoldco/promptui"
	ptime "github.com/yaa110/go-persian-calendar"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"ravro_dcrpt/core"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	rvrVersion    = "v1.0.2"
	publicMessage = "شرح داده نشده است."
	noMsg         = "ثبت نشد"
	PathDir       = "encrypt"
	keyDir        = "key"
	DecryptPath   = "decrypt"
	ExtRavro      = "*.ravro"
	ExtZip        = "*.zip"
)

func main() {

	var (
		templatePath string
		outputPath   string
		keyFixPath   string
		outFixpath   string
		curretnPath  string
		status       bool
		files        []fs.FileInfo
	)
	lstDir := []string{"encrypt", "decrypt", "key"}
	if runtime.GOOS == "windows" {
		templatePath = "template\\sample.html"
		outputPath = "decrypt\\reports.pdf"
		keyFixPath = "key/%s"
		outFixpath = "decrypt"
	} else {
		templatePath = "template/sample.html"
		outputPath = "decrypt/reports.pdf"
		keyFixPath = "key/%s"
		outFixpath = "decrypt"
	}
	init := flag.Bool("init", false, "Create encrypt/decrypt/key directory: ./ravro_dcrpt -init")
	inputDir := flag.String("in", "in", "Input directory of report encrypt file, Ex: ./ravro_dcrpt -in=/home/path")
	outputDir := flag.String("out", "out", "Output directory for decrypt report file,Ex: ./ravro_dcrpt -out=/home/path/")
	key := flag.String("key", "key", "Store key name Ex: ./ravro_dcrpt -key=/home/key.private")

	update := flag.Bool("update", false, "Update ravro decryptor")
	format := flag.Bool("json", false, "Convert report to json")
	logs := flag.Bool("log", false, "Store Error logs in log.txt")
	fmt.Println(">> Help : ravro_dcrpt --help")
	fmt.Println(">> Current Version : ravro_dcrpt/1.0.2")
	fmt.Println(">> Github : https://github.com/ravro-ir/ravro_dcrp")
	fmt.Println(">> Issue : https://github.com/ravro-ir/ravro_dcrp/issues")
	fmt.Println(">> Author : Ravro Development Team (RDT) \n\n")

	flag.Parse()

	ver, err := utils.LatestVersion()
	if err != nil {
		LogCheck(*logs, err)
	} else {
		newVer := fmt.Sprintf("ravro_dcrpt/%s", ver)
		if rvrVersion != ver {
			fmt.Println(fmt.Sprintf("\n New version (%s) released, "+
				"Please use command : ./ravro_dcrpt -update \n", newVer))
		}
	}
	if *update {
		fmt.Println("[++++] Downloading latest version from Github")
		fileName, verTag, err := utils.HttpGet()
		fmt.Println(fmt.Sprintf("[++++] Updating from [%s] -> [%s]", rvrVersion, verTag))
		if err != nil {
			LogCheck(*logs, err)
			fmt.Println("[----] Unable to download latest file, Maybe your internet connection is bad," +
				"or please usage : `./ravro_dcrpt -log` for see monitoring error message")
		}
		fmt.Println(fmt.Sprintf("[++++] The latest version Ravro_dcrpt downloaded - [%s]", fileName))
		// Extract zip file
		path, err := os.Getwd()
		newPathZip := filepath.Join(path, fileName)
		err = utils.Unzip(newPathZip, path)
		if err != nil {
			LogCheck(*logs, err)
			fmt.Println("[----] Error : Unable to extract zip file.")
		}
		return
	}
	if *init {
		utils.AddDir(DecryptPath)
		utils.AddDir(PathDir)
		utils.AddDir(keyDir)
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
	if *key != keyDir {
		status = true

		keyFixPath = *key
	}
	var path string
	if curretnPath == "" {
		path, err = utils.Projectpath()
		if err != nil {
			LogCheck(*logs, err)
			fmt.Println(err)
			return
		}
	} else {
		path = curretnPath
	}

	if *key == keyDir {
		keyPath := filepath.Join(path, keyDir)
		files, err = ioutil.ReadDir(keyPath)
		if err != nil {
			LogCheck(*logs, err)
			fmt.Println(err)
			return
		}
		if len(files) == 1 {
			keyFixPath = fmt.Sprintf(keyFixPath, files[0].Name())
		} else {
			var lst []string
			for _, value := range files {
				lst = append(lst, value.Name())
			}
			prompt := promptui.Select{
				Label: "Please choose a key",
				Items: lst,
			}
			_, result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			keyFixPath = fmt.Sprintf(keyFixPath, result)
		}
	}

	lstReport, err := utils.ReportFiles(path, ExtRavro)
	if err != nil {
		LogCheck(*logs, err)
		fmt.Println(err)
		return
	}
	CurrPath, _ := os.Getwd()
	var ll []string
	for _, value := range lstReport {
		reportPath := filepath.Join(CurrPath, PathDir, utils.GetReportID(value))
		ll = append(ll, reportPath)
	}
	zipFile, err := utils.ReportFiles(path, ExtZip)
	if err != nil {
		return
	}
	var lstZipFilepath []string
	if len(ll) >= 1 {
		lstZipFilepath = append(lstZipFilepath, utils.Unique(ll)...)
	} else {

	}
	var extractPath string
	for _, value := range zipFile {

		if curretnPath == "" {
			extractPath = filepath.Join(CurrPath, PathDir, utils.GetReportID(value))
		} else {
			extractPath = curretnPath + utils.GetReportID(value)
		}
		err := utils.Unzip(value, extractPath)
		if err != nil {
			fmt.Println("[----] Error : Unable to extract zip file.")
			return
		}
		lstZipFilepath = append(lstZipFilepath, extractPath)
	}
	r := utils.NewRequestPdf("")
	pt := ptime.Now()
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "ID", "Hunter", "Status")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "-------", "-------", "-------")

	checkList := make(map[string]string)
	for _, zipdata := range utils.Unique(lstZipFilepath) {
		curretnPath = zipdata
		fmt.Println(fmt.Sprintf("[++++] Starting for decrypting report ID [%s] . . . ", utils.GetReportID(zipdata)))
		report, err := core.DcrptReport(curretnPath, keyFixPath, outFixpath, status)
		if err != nil {
			LogCheck(*logs, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed"
			fmt.Println("[----] Error : Unable to decrypt files, We think your key is invalid. Please use : ./ravro_dcrpt -log")
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			continue
		}
		if report.Title == "" {
			LogCheck(*logs, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed"
			fmt.Println("[----] The input file for decryption is not correct.")
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			continue
		}
		fmt.Println("[++++] Starting for decrypting Judgment . . . ")
		judge, err := core.DcrptJudgment(curretnPath, keyFixPath, outFixpath, status)
		if err != nil {
			LogCheck(*logs, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed"
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			continue
		}
		fmt.Println("[++++] Starting for decrypting Amendment . . . ")
		amendment, err := core.DcrptAmendment(curretnPath, keyFixPath, outFixpath)
		if err != nil {
			LogCheck(*logs, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed"
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			continue
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
		if *format {
			file, _ := json.MarshalIndent(pdf.Judge, "", " ")
			fileJurorsonPath := filepath.Join("decrypt", "juror.json")
			_ = ioutil.WriteFile(fileJurorsonPath, file, 0644)
			reportd, _ := json.MarshalIndent(pdf.Report, "", " ")
			fileRepoPath := filepath.Join("decrypt", "repo.json")
			_ = ioutil.WriteFile(fileRepoPath, reportd, 0644)
			amendments, _ := json.MarshalIndent(pdf.Amendment, "", " ")
			fileMorePath := filepath.Join("decrypt", "moreinfo.json")
			_ = ioutil.WriteFile(fileMorePath, amendments, 0644)
		}
		md := []byte(pdf.Report.Description)
		templateData := TemplateStruct(md, pdf, dateFrom, dateTo, moreInfo, report)
		if err := r.ParseTemplate(templatePath, templateData); err == nil {
			s := spinner.New(spinner.CharSets[4], 100*time.Millisecond) // Build our new spinner
			s.Start()
			s.Color("yellow")
			s.Prefix = "[++++] Starting report to pdf "
			_, err = r.GeneratePDF(outputPath)
			if err != nil {
				LogCheck(*logs, err)
				checkList["id"] = report.Slug
				checkList["hunter"] = report.HunterUsername
				checkList["status"] = "Failed"
				fmt.Println("[----] failed to remove html template,")
				fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
				continue
			}
			err := os.RemoveAll("template")
			if err != nil {
				LogCheck(*logs, err)
				checkList["id"] = report.Slug
				checkList["hunter"] = report.HunterUsername
				checkList["status"] = "Failed"
				fmt.Println("[----] failed to remove html template,")
				fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
				continue
			}
			fmt.Println("\n[++++] PDF generated successfully")
			fmt.Println("\n")
			err = utils.ChangeDirName(report.Slug, outFixpath)
			if err != nil {
				LogCheck(*logs, err)
				checkList["id"] = report.Slug
				checkList["hunter"] = report.HunterUsername
				checkList["status"] = "Failed"
				fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
				continue
			}
			s.Stop()
		} else {
			LogCheck(*logs, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed"
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			continue
		}
		checkList["id"] = report.Slug
		checkList["hunter"] = report.HunterUsername
		checkList["status"] = "Successfully"
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
		fmt.Println("\n")
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
	var valuation string
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
	if report.ReportInfo.Details.CurrentStatus == "" {
		report.ReportInfo.Details.CurrentStatus = noMsg
	}
	score, _ := strconv.ParseFloat(report.ReportInfo.Details.Cvss.Hunter.Score, 32)

	if report.Urls == "" {
		report.Urls = noMsg
	}
	if report.ReportInfo.Details.Target == "" {
		report.ReportInfo.Details.Target = noMsg
	}
	if score > 1 && score < 4 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "Low"
	}
	if score >= 4 && score < 7 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "High"
	}
	if score >= 7 && score < 9 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "Medium"
	}
	if score >= 9 && score < 10 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "Critical"
	}
	if report.ReportInfo.Details.CurrentStatus == "بررسی اولیه" {
		valuation = fmt.Sprintf("<a href=\"https://www.ravro.ir/fa/report/%s/valuation\" ><img class=\"img-center2\" width=\"200px\" src=\" data:image/png;base64,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\" /></a>\n", report.Slug)
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
		UrlTarget       string
		Valuation       string
	}{
		Title:           pdf.Report.Title,
		PoC:             string(output),
		CVSSJudge:       pdf.Judge.Cvss.Value,
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
		ScoreJudge:      pdf.Judge.Cvss.Rating,
		ScoreHunter:     report.ReportInfo.Details.Cvss.Hunter.Score,
		CVSSHunter:      report.ReportInfo.Details.Cvss.Hunter.Vector,
		RangeDate:       "از" + report.DateFrom + "  تا " + report.DateTo,
		Targets:         report.ReportInfo.Details.Target,
		Status:          report.ReportInfo.Details.CurrentStatus,
		UrlTarget:       report.Urls,
		Valuation:       valuation,
	}
	return templateData
}

func LogCheck(status bool, Msg error) {
	if status {
		utils.Logger(Msg)
	}
}
