//go:build linux || darwin
// +build linux darwin

package pdfgen

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	wk "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	ptime "github.com/yaa110/go-persian-calendar"

	"ravro_dcrpt/internal/ports"
	"ravro_dcrpt/pkg/models"
)

// WKHTMLToPDFGenerator implements PDFGenerator using wkhtmltopdf
type WKHTMLToPDFGenerator struct {
	body string
}

// NewWKHTMLToPDFGenerator creates a new wkhtmltopdf-based PDF generator
func NewWKHTMLToPDFGenerator() ports.PDFGenerator {
	return &WKHTMLToPDFGenerator{}
}

// SetTemplate sets the template path (not used in this implementation)
func (g *WKHTMLToPDFGenerator) SetTemplate(templatePath string) error {
	return nil
}

// GenerateReport generates a PDF report using HTML template
func (g *WKHTMLToPDFGenerator) GenerateReport(report *models.Report, judgment *models.Judgment, amendments []models.Amendment, outputPath string) error {
	// Create template data
	templateData := createTemplateData(report, judgment, amendments)

	// Parse HTML template
	tmpl, err := template.New("report").Parse(HTMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute template
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, templateData); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.body = buf.String()

	// Clean up HTML entities
	g.body = strings.ReplaceAll(g.body, "&#34;&lt;", "<")
	g.body = strings.ReplaceAll(g.body, "&gt;", ">")
	g.body = strings.ReplaceAll(g.body, "&lt;", "<")
	g.body = strings.ReplaceAll(g.body, "&#34;", "")

	// Generate PDF
	return g.generatePDF(outputPath)
}

// generatePDF generates the PDF file from the parsed HTML
func (g *WKHTMLToPDFGenerator) generatePDF(pdfPath string) error {
	// Validate HTML content
	if len(strings.TrimSpace(g.body)) == 0 {
		return fmt.Errorf("HTML content is empty")
	}

	// Create temporary directory
	tmpPath := "template/"
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		if err := os.Mkdir(tmpPath, 0755); err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	// Write HTML to temp file
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	htmlFileName := tmpPath + timestamp + ".html"
	if err := os.WriteFile(htmlFileName, []byte(g.body), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	// Remove HTML file after PDF generation
	defer func() {
		os.Remove(htmlFileName)
		// Try to remove template directory if empty
		os.Remove(tmpPath)
	}()

	// Create PDF generator
	pdfg, err := wk.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to create PDF generator: %w", err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wk.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	// Enable external resources (for fonts, images, etc.)
	pdfg.NoCollate.Set(false)

	// Create page from HTML file
	page := wk.NewPage(htmlFileName)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Enable loading external resources (fonts from Google, etc.)
	page.EnableLocalFileAccess.Set(true)
	page.LoadErrorHandling.Set("ignore")
	page.LoadMediaErrorHandling.Set("ignore")

	pdfg.AddPage(page)

	// Create PDF
	if err := pdfg.Create(); err != nil {
		return fmt.Errorf("failed to create PDF: %w", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(pdfPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write PDF to file
	if err := pdfg.WriteFile(pdfPath); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

// GenerateFromHTML generates a PDF from HTML content (not implemented)
func (g *WKHTMLToPDFGenerator) GenerateFromHTML(htmlContent string, outputPath string) error {
	g.body = htmlContent
	return g.generatePDF(outputPath)
}

// convertToShamsi converts Gregorian date string (YYYY-MM-DD) to Persian date
func convertToShamsi(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Parse the date string (format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)
	parts := strings.Split(dateStr, " ")
	datePart := parts[0]
	dateFields := strings.Split(datePart, "-")

	if len(dateFields) != 3 {
		return dateStr // Return original if format is unexpected
	}

	year, err1 := strconv.Atoi(dateFields[0])
	month, err2 := strconv.Atoi(dateFields[1])
	day, err3 := strconv.Atoi(dateFields[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return dateStr // Return original if conversion fails
	}

	// Convert to Persian date
	t := time.Date(year, time.Month(month), day, 12, 0, 0, 0, ptime.Iran())
	pt := ptime.New(t)
	return pt.Format("yyyy/MM/dd")
}

// formatAmount formats amount with thousand separators
func formatAmount(amount int) string {
	if amount == 0 {
		return "در حال بررسی"
	}

	// Convert to string
	str := strconv.Itoa(amount)

	// Add thousand separators
	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}

	return result.String() + " ریال"
}

// markdownToHTML converts Markdown text to HTML
func markdownToHTML(mdText string) string {
	if strings.TrimSpace(mdText) == "" {
		return ""
	}

	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Parse markdown
	doc := p.Parse([]byte(mdText))

	// Create HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Render to HTML
	htmlBytes := markdown.Render(doc, renderer)

	return string(htmlBytes)
}

// createTemplateData creates template data structure from report and judgment
func createTemplateData(report *models.Report, judgment *models.Judgment, amendments []models.Amendment) interface{} {
	// Extract judge names
	judgeNames := ""
	if len(report.ReportInfo.Details.Judges) > 0 {
		for i, judge := range report.ReportInfo.Details.Judges {
			if i > 0 {
				judgeNames += ", "
			}
			judgeNames += judge.Name
		}
	}

	// Extract attachment list as HTML table
	attachmentList := ""
	if len(report.ReportInfo.Details.Attachments) > 0 {
		attachmentList = `<table class="data-table no-break"><thead><tr><th>ردیف</th><th>نام فایل</th><th>نوع</th></tr></thead><tbody>`
		for i, att := range report.ReportInfo.Details.Attachments {
			fileType := "فایل"
			if strings.Contains(att.Filename, ".png") || strings.Contains(att.Filename, ".jpg") || strings.Contains(att.Filename, ".jpeg") {
				fileType = "تصویر"
			} else if strings.Contains(att.Filename, ".pdf") {
				fileType = "PDF"
			}
			attachmentList += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", i+1, att.Filename, fileType)
		}
		attachmentList += `</tbody></table>`
	}

	// Extract tags info
	moreInfo := ""
	linkMoreInfo := ""
	if len(report.ReportInfo.Tags) > 0 {
		for _, tag := range report.ReportInfo.Tags {
			if tag.InfoDescription != "" {
				moreInfo += tag.InfoDescription + "\n\n"
			}
			if tag.InfoTitle != "" && tag.InfoSolution != "" {
				moreInfo += "### " + tag.InfoTitle + "\n\n"
				moreInfo += tag.InfoSolution + "\n\n"
			}
			if tag.InfoMore != "" {
				linkMoreInfo += tag.InfoMore + "\n"
			}
		}
	}

	// Extract amendment info (نیازمندی اطلاعات بیشتر از طرف شرکت)
	if len(amendments) > 0 {
		if moreInfo != "" {
			moreInfo += "\n\n---\n\n"
		}

		// Amendment info without title
		for i, amendment := range amendments {
			content := amendment.Description
			if content == "" {
				content = amendment.Content // Fallback to legacy field
			}
			if content != "" {
				if i > 0 {
					moreInfo += "\n\n"
				}
				moreInfo += content
			}
		}
	}

	// Convert markdown to HTML for info fields
	if moreInfo != "" {
		moreInfo = markdownToHTML(moreInfo)
	}
	if linkMoreInfo != "" {
		linkMoreInfo = markdownToHTML(linkMoreInfo)
	}

	// Prepare judgment data
	judgeUser := judgeNames
	amount := 0
	cvssJudge := ""
	scoreJudge := ""
	judgeInfo := ""

	if judgment != nil {
		amount = judgment.Reward
		cvssJudge = judgment.Cvss.Value
		scoreJudge = judgment.Cvss.Rating
		judgeInfo = judgment.Description

		// Also check vulnerability info
		if judgment.Vulnerability.Define != "" {
			if judgeInfo != "" {
				judgeInfo += "\n\n"
			}
			judgeInfo += "**📝 تعریف:**\n\n" + judgment.Vulnerability.Define
		}
		if judgment.Vulnerability.Fix != "" {
			if judgeInfo != "" {
				judgeInfo += "\n\n"
			}
			judgeInfo += "**🔧 راه حل:**\n\n" + judgment.Vulnerability.Fix
		}
	}

	// Also use judge CVSS from report if available
	if cvssJudge == "" && report.ReportInfo.Details.Cvss.Judge.Vector != "" {
		cvssJudge = report.ReportInfo.Details.Cvss.Judge.Vector
		scoreJudge = report.ReportInfo.Details.Cvss.Judge.Score
	}

	// Convert markdown to HTML for judgment info
	if judgeInfo != "" {
		judgeInfo = markdownToHTML(judgeInfo)
	}

	// Convert markdown to HTML for report fields
	scenario := markdownToHTML(report.Scenario)
	description := markdownToHTML(report.Description)

	return struct {
		Title           string
		Hunter          string
		ReportID        string
		CompanyUserName string
		DateFrom        string
		Status          string
		Targets         string
		CVSSHunter      string
		ScoreHunter     string
		Ips             string
		RangeDate       string
		Scenario        template.HTML
		UrlTarget       string
		PoC             template.HTML
		MoreInfo        template.HTML
		JudgeUser       string
		Amount          string
		CVSSJudge       string
		ScoreJudge      string
		DateTo          string
		JudgeInfo       template.HTML
		LinkMoreInfo    template.HTML
		Attachment      template.HTML
		RavroVer        string
		Valuation       string
	}{
		Title:           report.Title,
		Hunter:          report.HunterUsername,
		ReportID:        report.Slug,
		CompanyUserName: report.CompanyUsername,
		DateFrom:        convertToShamsi(report.SubmissionDate),
		Status:          report.ReportInfo.Details.CurrentStatus,
		Targets:         report.ReportInfo.Details.Target,
		CVSSHunter:      report.ReportInfo.Details.Cvss.Hunter.Vector,
		ScoreHunter:     report.ReportInfo.Details.Cvss.Hunter.Score,
		Ips:             report.Ips,
		RangeDate:       fmt.Sprintf("از %s تا %s", convertToShamsi(report.DateFrom), convertToShamsi(report.DateTo)),
		Scenario:        template.HTML(scenario),
		UrlTarget:       report.Urls,
		PoC:             template.HTML(description),
		MoreInfo:        template.HTML(moreInfo),
		JudgeUser:       judgeUser,
		Amount:          formatAmount(amount),
		CVSSJudge:       cvssJudge,
		ScoreJudge:      scoreJudge,
		DateTo:          convertToShamsi(report.DateTo),
		JudgeInfo:       template.HTML(judgeInfo),
		LinkMoreInfo:    template.HTML(linkMoreInfo),
		Attachment:      template.HTML(attachmentList),
		RavroVer:        "v2.0.0",
		Valuation:       "",
	}
}
