package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gomarkdown/markdown"
	"github.com/manifoldco/promptui"
	ptime "github.com/yaa110/go-persian-calendar"
	"ravro_dcrpt/core"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
)

const (
	rvrVersion       = "v1.0.4"
	publicMessage    = "شرح داده نشده است."
	noMsg            = "ثبت نشد"
	PathDir          = "encrypt"
	keyDir           = "key"
	DecryptPath      = "decrypt"
	ExtRavro         = "*.ravro"
	ExtZip           = "*.zip"
	defaultInputDir  = "in"
	defaultOutputDir = "decrypt"
	defaultKeyDir    = "key"
)

type Config struct {
	TemplatePath string
	OutputPath   string
	KeyPath      string
	OutputDir    string
	InputDir     string
	Init         bool
	Update       bool
	JSONFormat   bool
	Logs         bool
}

func main() {
	// Parse command-line flags
	config := parseFlags()

	// Display help information
	displayHelp()

	// Check for new version
	err := handleVersionCheck()
	if err != nil {
		logError(config.Logs, err)
	}

	// Handle update if requested
	if config.Update {
		err := handleUpdate(config.Logs)
		if err != nil {
			logError(config.Logs, err)
			fmt.Println("[----] Unable to check for latest version. Please use: ./ravro_dcrpt -log for error details")
		}
		return
	}

	// Initialize directories if requested
	if config.Init {
		initializeDirectories()
		return
	}

	// Validate directories
	err = validateDirectories()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Process paths
	err = processPaths(&config)
	if err != nil {
		logError(config.Logs, err)
		fmt.Println("[----] Error occurred. More Info: ./ravro_dcrpt -log")
		return
	}

	// Select key if needed
	err = selectKeyIfNeeded(&config)
	if err != nil {
		logError(config.Logs, err)
		fmt.Println("[----] Error occurred. More Info: ./ravro_dcrpt -log")
		return
	}

	// Collect report paths
	reportPaths, err := collectReportPaths(config.InputDir, config.Logs)
	if err != nil {
		fmt.Println("[----] Error occurred. More Info: ./ravro_dcrpt -log")
		return
	}

	// Initialize table writer for output
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "ID", "Hunter", "Status")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "-------", "-------", "-------")

	// Process each report
	for _, reportPath := range reportPaths {
		if !strings.Contains(reportPath, "/encrypt/") {
			fmt.Printf("[----] Skipping path without /encrypt/: %s\n", reportPath)
			continue
		}
		processReport(reportPath, config.KeyPath, config.OutputDir, config.JSONFormat, config.Logs, w)
	}
}

// parseFlags parses command-line flags and returns a Config struct
func parseFlags() Config {
	config := Config{
		TemplatePath: filepath.Join("template", "sample.html"),
		OutputPath:   filepath.Join("decrypt", "reports.pdf"),
		OutputDir:    "decrypt",
	}

	init := flag.Bool("init", false, "Create encrypt/decrypt/key directories")
	inputDir := flag.String("in", defaultInputDir, "Input directory of encrypted report files")
	outputDir := flag.String("out", defaultOutputDir, "Output directory for decrypted reports")
	key := flag.String("key", defaultKeyDir, "Path to key file")
	update := flag.Bool("update", false, "Update ravro decryptor")
	myJson := flag.Bool("json", false, "Convert report to JSON")
	logs := flag.Bool("log", false, "Store error logs in log.txt")

	flag.Parse()

	config.Init = *init
	config.Update = *update
	config.JSONFormat = *myJson
	config.Logs = *logs
	config.InputDir = *inputDir
	config.OutputDir = *outputDir
	config.KeyPath = *key

	return config
}

// displayHelp prints help information to the console
func displayHelp() {
	fmt.Println(">> Help : ravro_dcrpt --help")
	fmt.Printf(">> Current Version : ravro_dcrpt/%s\n", rvrVersion)
	fmt.Println(">> Github : https://github.com/ravro-ir/ravro_dcrp")
	fmt.Println(">> Issue : https://github.com/ravro-ir/ravro_dcrp/issues")
	fmt.Println(">> Author : Ravro Development Team (RDT)\n")
}

// handleVersionCheck checks for a new version and returns an error if any
func handleVersionCheck() error {
	ver, err := utils.LatestVersion()
	if err != nil {
		return err
	}
	if rvrVersion != ver {
		fmt.Printf("\nNew version (ravro_dcrpt/%s) released, use: ./ravro_dcrpt -update\n", ver)
	}
	return nil
}

// handleUpdate updates the application to the latest version
func handleUpdate(logs bool) error {
	fmt.Println("[++++] Checking for latest version from GitHub")
	latestVer, err := utils.LatestVersion()
	if err != nil {
		return err
	}

	if rvrVersion == latestVer {
		fmt.Println("[++++] You have the latest version of ravro_dcrpt")
		return nil
	}

	fileName, verTag, err := utils.HttpGet()
	if err != nil {
		return err
	}

	fmt.Printf("[++++] Updating from [%s] -> [%s]\n", rvrVersion, verTag)

	path, _ := os.Getwd()
	if err := utils.Unzip(filepath.Join(path, fileName), path); err != nil {
		return err
	}

	fmt.Printf("[++++] Latest version downloaded - [%s]\n", fileName)
	return nil
}

// initializeDirectories creates necessary directories
func initializeDirectories() {
	dirs := []string{PathDir, DecryptPath, keyDir}
	for _, dir := range dirs {
		utils.AddDir(dir)
	}
	fmt.Println("[++] Created directories: decrypt, encrypt, key")
}

// validateDirectories ensures all required directories exist
func validateDirectories() error {
	dirs := []string{PathDir, DecryptPath, keyDir}
	for _, dir := range dirs {
		if utils.CheckDir(dir) {
			return fmt.Errorf("[---] Required directories missing. Usage: ravro_dcrpt -init")
		}
	}
	return nil
}

// processPaths updates config paths based on user input
func processPaths(config *Config) error {
	if config.InputDir != defaultInputDir {
		// Keep user provided path
	} else {
		path, err := utils.Projectpath()
		if err != nil {
			return err
		}
		config.InputDir = path
	}

	if config.OutputDir != defaultOutputDir {
		config.OutputPath = filepath.Join(config.OutputDir, "reports.pdf")
	}

	return nil
}

// selectKeyIfNeeded prompts for key selection if multiple keys found
func selectKeyIfNeeded(config *Config) error {
	if config.KeyPath == defaultKeyDir {
		currentPath, _ := os.Getwd()
		keyPath := filepath.Join(currentPath, defaultKeyDir)
		files, err := ioutil.ReadDir(keyPath)
		if err != nil {
			return err
		}

		switch len(files) {
		case 0:
			return fmt.Errorf("no key files found in key directory")
		case 1:
			config.KeyPath = filepath.Join(keyDir, files[0].Name())
		default:
			keys := make([]string, len(files))
			for i, f := range files {
				keys[i] = f.Name()
			}
			prompt := promptui.Select{
				Label: "Please choose a key",
				Items: keys,
			}
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			config.KeyPath = filepath.Join(keyDir, result)
		}
	}
	return nil
}

// logError logs an error if logging is enabled
func logError(logEnabled bool, err error) {
	if logEnabled {
		utils.Logger(err)
	}
}

// arrayStringToInt converts a string to an integer
func arrayStringToInt(date string) int {
	newData, err := strconv.Atoi(date)
	if err != nil {
		log.Fatal(err)
	}
	return newData
}

// miladiToShamsi converts Gregorian date to Persian date
func miladiToShamsi(year, month, day int) string {
	var t = time.Date(year, time.Month(month), day, 12, 1, 1, 0, ptime.Iran())
	pt := ptime.New(t)
	shamsiDate := pt.Format("yyyy/MM/dd")
	return shamsiDate
}

// attachmentFiles formats attachment information for reports
func attachmentFiles(info entity.InfoReport) string {
	var attach string
	if len(info.Details.Attachments) == 0 {
		return "<tr>\n<td>فایلی توسط شکارچی پیوست نشده است</td>\n</tr>"
	}
	for _, content := range info.Details.Attachments {
		attach += fmt.Sprintf("<tr>\n<td>%s</td>\n</tr>", content.Filename)
	}
	return attach
}

// judgmentUser formats judgment user information
func judgmentUser(info entity.InfoReport) string {
	var judge string
	for _, content := range info.Details.Judges {
		judge += content.Name + "  "
	}
	return judge
}

// conString formats additional information for reports
func conString(info entity.InfoReport) string {
	var infoMore string
	for _, content := range info.Tags {
		infoMore += " عنوان تگ : " + content.InfoTitle + "<br />"
		infoMore += " توضیحات گزارش : " + "<br />" + content.InfoDescription + "<br />"
		infoMore += "راه حل : " + "<br />" + content.InfoSolution + "<br />"
		infoMore += "اطلاعات بیشتر : " + "<br />" + content.InfoMore + "<br />"
		infoMore += "<hr>"
	}
	return infoMore
}

// validateReport validates report data and returns date and output path
func validateReport(report entity.Report, outputPath string, pdf entity.Pdf) (string, string) {
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
		log.Fatalln("[----] Error: The submit date is empty, we think your report path is incorrect, (Valid Path: (encrypt/ir2022-01-10-0001))")
	}

	dateReport := strings.Split(dateSubmited[0], "-")
	dateFrom := miladiToShamsi(arrayStringToInt(dateReport[0]), arrayStringToInt(dateReport[1]), arrayStringToInt(dateReport[2]))
	return dateFrom, outputPath
}

// prepareTemplateData prepares template data for PDF generation
func prepareTemplateData(md []byte, senario []byte, pdf entity.Pdf, dateFrom, dateTo, moreInfo string, report entity.Report) interface{} {
	var valuation string
	output := markdown.ToHTML(md, nil, nil)
	infoTrain := conString(report.ReportInfo)
	pdf.Report.Scenario = string(markdown.ToHTML(senario, nil, nil))

	// Process date formats
	dateOf := strings.Split(report.DateFrom, "-")
	if len(dateOf) == 0 {
		log.Fatalln("[----] Error: The datefrom is invalid")
	}
	report.DateFrom = miladiToShamsi(arrayStringToInt(dateOf[0]), arrayStringToInt(dateOf[1]), arrayStringToInt(dateOf[2]))

	dateTwo := strings.Split(report.DateTo, "-")
	if len(dateTwo) == 0 {
		log.Fatalln("[----] Error: The datefrom is invalid")
	}
	report.DateTo = miladiToShamsi(arrayStringToInt(dateTwo[0]), arrayStringToInt(dateTwo[1]), arrayStringToInt(dateTwo[2]))

	// Set default values if needed
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

	// Process CVSS score
	score, _ := strconv.ParseFloat(report.ReportInfo.Details.Cvss.Hunter.Score, 32)

	if report.Urls == "" {
		report.Urls = noMsg
	}
	if report.ReportInfo.Details.Target == "" {
		report.ReportInfo.Details.Target = noMsg
	}

	// Map numerical scores to text categories
	if score > 1 && score < 4 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "Low"
	}
	if score >= 4 && score < 7 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "Medium"
	}
	if score >= 7 && score < 9 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "High"
	}
	if score >= 9 && score < 10 {
		report.ReportInfo.Details.Cvss.Hunter.Score = "Critical"
	}

	if report.ReportInfo.Details.CurrentStatus == "بررسی اولیه" {
		valuation = fmt.Sprintf("<a href=\"https://www.ravro.ir/fa/report/%s/valuation\" ", report.Slug)
	}

	// Create template data structure
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
		Attachment:      attachmentFiles(report.ReportInfo),
		Scenario:        pdf.Report.Scenario,
		LinkMoreInfo:    infoTrain,
		RavroVer:        rvrVersion,
		JudgeUser:       judgmentUser(report.ReportInfo),
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

// processReport handles decryption and processing of a single report
func processReport(reportPath, keyPath, outputPath string, jsonFormat, logEnabled bool, w *tabwriter.Writer) {
	fmt.Printf("[++++] Starting for decrypting report ID [%s] . . . \n", utils.GetReportID(reportPath))

	checkList := make(map[string]string)

	// Decrypt report
	report, err := core.DcrptReport(reportPath, keyPath, outputPath, true)
	if err != nil {
		logError(logEnabled, err)
		checkList["id"] = report.Slug
		checkList["hunter"] = report.HunterUsername
		checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
		fmt.Println("[----] Error: Unable to decrypt files, We think your key is invalid. Please use: ./ravro_dcrpt -log")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
		return
	}

	if report.Title == "" {
		logError(logEnabled, fmt.Errorf("empty report title"))
		checkList["id"] = report.Slug
		checkList["hunter"] = report.HunterUsername
		checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
		fmt.Println("[----] The input file for decryption is not correct. More Info: ./ravro_dcrpt -log")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
		return
	}

	// Decrypt judgment
	fmt.Println("[++++] Starting for decrypting Judgment . . . ")
	judge, err := core.DcrptJudgment(reportPath, keyPath, outputPath, true)
	if err != nil {
		logError(logEnabled, err)
		checkList["id"] = report.Slug
		checkList["hunter"] = report.HunterUsername
		checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
		return
	}

	// Decrypt amendment
	fmt.Println("[++++] Starting for decrypting Amendment . . . ")
	amendment, err := core.DcrptAmendment(reportPath, keyPath, outputPath)
	if err != nil {
		logError(logEnabled, err)
		checkList["id"] = report.Slug
		checkList["hunter"] = report.HunterUsername
		checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
		return
	}

	// Prepare for PDF generation
	templatePath := filepath.Join("template", "sample.html")
	utils.AddDir("template")
	utils.HtmlTemplate(templatePath)

	moreInfo := strings.Join(amendment[:], ",")
	if moreInfo == "" {
		moreInfo = publicMessage
	}
	if report.Reproduce == "" {
		report.Reproduce = publicMessage
	}

	// Get current Persian date
	pt := ptime.Now()
	dateTo := strconv.Itoa(pt.Year()) + "/" + strconv.Itoa(int(pt.Month())) + "/" + strconv.Itoa(pt.Day())

	pdf := entity.Pdf{Judge: judge, Report: report}

	// Validate report data
	dateFrom, finalOutputPath := validateReport(report, filepath.Join(outputPath, "reports.pdf"), pdf)

	// Generate JSON if requested
	if jsonFormat {
		generateJSONFiles(pdf, logEnabled)
	}

	// Generate PDF
	r := utils.NewRequestPdf("")
	md := []byte(pdf.Report.Description)
	scenario := []byte(pdf.Report.Scenario)
	templateData := prepareTemplateData(md, scenario, pdf, dateFrom, dateTo, moreInfo, report)

	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		s := spinner.New(spinner.CharSets[4], 100*time.Millisecond)
		s.Start()
		s.Color("yellow")
		s.Prefix = "[++++] Starting report to pdf "

		_, err = r.GeneratePDF(finalOutputPath)
		if err != nil {
			logError(logEnabled, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			s.Stop()
			return
		}

		err := os.RemoveAll("template")
		if err != nil {
			logError(logEnabled, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			s.Stop()
			return
		}

		fmt.Println("\n[++++] PDF generated successfully")
		fmt.Println("\n")

		err = utils.ChangeDirName(report.Slug, outputPath)
		if err != nil {
			logError(logEnabled, err)
			checkList["id"] = report.Slug
			checkList["hunter"] = report.HunterUsername
			checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
			s.Stop()
			return
		}
		s.Stop()
	} else {
		logError(logEnabled, err)
		checkList["id"] = report.Slug
		checkList["hunter"] = report.HunterUsername
		checkList["status"] = "Failed, More Info: ./ravro_dcrpt -log"
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
		return
	}

	checkList["id"] = report.Slug
	checkList["hunter"] = report.HunterUsername
	checkList["status"] = "Successfully"
	fmt.Fprintf(w, "\n %s\t%s\t%s\t", checkList["id"], checkList["hunter"], checkList["status"])
	fmt.Println("\n")
}

// generateJSONFiles generates JSON files from PDF data
func generateJSONFiles(pdf entity.Pdf, logEnabled bool) {
	// Generate JSON for judge data
	fileJudge, err := json.MarshalIndent(pdf.Judge, "", " ")
	if err != nil {
		logError(logEnabled, err)
		fmt.Println("[----] We have a error, More Info: ./ravro_dcrpt -log")
		return
	}
	fileJurorsonPath := filepath.Join("decrypt", "juror.json")
	err = ioutil.WriteFile(fileJurorsonPath, fileJudge, 0644)
	if err != nil {
		logError(logEnabled, err)
		fmt.Println("[----] We have a error, More Info: ./ravro_dcrpt -log")
		return
	}

	// Generate JSON for report data
	reportData, err := json.MarshalIndent(pdf.Report, "", " ")
	if err != nil {
		logError(logEnabled, err)
		fmt.Println("[----] We have a error, More Info: ./ravro_dcrpt -log")
		return
	}
	fileRepoPath := filepath.Join("decrypt", "repo.json")
	err = ioutil.WriteFile(fileRepoPath, reportData, 0644)
	if err != nil {
		logError(logEnabled, err)
		fmt.Println("[----] We have a error, More Info: ./ravro_dcrpt -log")
		return
	}

	// Generate JSON for amendment data
	amendments, err := json.MarshalIndent(pdf.Amendment, "", " ")
	if err != nil {
		logError(logEnabled, err)
		fmt.Println("[----] We have a error, More Info: ./ravro_dcrpt -log")
		return
	}
	fileMorePath := filepath.Join("decrypt", "moreinfo.json")
	err = ioutil.WriteFile(fileMorePath, amendments, 0644)
	if err != nil {
		logError(logEnabled, err)
		fmt.Println("[----] We have a error, More Info: ./ravro_dcrpt -log")
	}
}

// collectReportPaths collects all report paths from zip and ravro files
func collectReportPaths(inputPath string, logEnabled bool) ([]string, error) {
	var reportPaths []string

	// Get current working directory
	currentPath, _ := os.Getwd()

	// Process .ravro files
	ravroFiles, err := utils.ReportFiles(inputPath, ExtRavro)
	if err != nil {
		logError(logEnabled, err)
		return nil, err
	}

	for _, file := range ravroFiles {
		reportID := utils.GetReportID(file)
		if reportID == "" {
			continue
		}
		reportPath := filepath.Join(currentPath, PathDir, reportID)
		reportPaths = append(reportPaths, reportPath)
	}

	// Process .zip files
	zipFiles, err := utils.ReportFiles(inputPath, ExtZip)
	if err != nil {
		logError(logEnabled, err)
		return nil, err
	}

	for _, file := range zipFiles {
		reportID := utils.GetReportID(file)
		if reportID == "" {
			continue
		}

		extractPath := filepath.Join(currentPath, PathDir, reportID)
		//if inputPath != "" {
		//	extractPath = filepath.Join(inputPath, reportID)
		//}

		err := utils.Unzip(file, extractPath)
		if err != nil {
			logError(logEnabled, err)
			fmt.Println("[----] Error: Unable to extract zip file, More Info: ./ravro_dcrpt -log")
			continue
		}

		reportPaths = append(reportPaths, extractPath)
	}

	return utils.Unique(reportPaths), nil
}
