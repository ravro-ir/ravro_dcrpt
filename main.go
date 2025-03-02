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
		normalizedPath := filepath.ToSlash(reportPath) // Converts Windows `\` to `/`
		if !strings.Contains(normalizedPath, "/encrypt/") {
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
		valuation = fmt.Sprintf("<a href=\"https://www.ravro.ir/fa/report/%s/valuation\" ><img class=\"img-center2\" width=\"200px\" src=\" data:image/png;base64,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\" /></a>\n", report.Slug)
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
