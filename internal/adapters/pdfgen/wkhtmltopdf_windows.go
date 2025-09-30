//go:build windows
// +build windows

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
func (g *WKHTMLToPDFGenerator) GenerateReport(report *models.Report, judgment *models.Judgment, outputPath string) error {
	// Create template data
	templateData := createTemplateData(report, judgment)

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
	defer os.Remove(htmlFileName)

	// Create PDF generator
	pdfg, err := wk.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to create PDF generator: %w", err)
	}

	// Note: wkhtmltopdf must be in PATH on Windows
	// Or set WKHTMLTOPDF_PATH environment variable

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wk.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	// Create page from HTML file
	page := wk.NewPage(htmlFileName)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)
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

// createTemplateData creates template data structure from report and judgment
func createTemplateData(report *models.Report, judgment *models.Judgment) interface{} {
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
		Scenario        string
		UrlTarget       string
		PoC             string
		MoreInfo        string
		JudgeUser       string
		Amount          int
		CVSSJudge       string
		ScoreJudge      string
		DateTo          string
		JudgeInfo       string
		LinkMoreInfo    string
		Attachment      string
		RavroVer        string
		Valuation       string
	}{
		Title:           report.Title,
		Hunter:          report.HunterUsername,
		ReportID:        report.Slug,
		CompanyUserName: report.CompanyUsername,
		DateFrom:        report.DateFrom,
		Status:          report.ReportInfo.Details.CurrentStatus,
		Targets:         report.ReportInfo.Details.Target,
		CVSSHunter:      report.ReportInfo.Details.Cvss.Hunter.Vector,
		ScoreHunter:     report.ReportInfo.Details.Cvss.Hunter.Score,
		Ips:             report.Ips,
		RangeDate:       fmt.Sprintf("از %s تا %s", report.DateFrom, report.DateTo),
		Scenario:        report.Scenario,
		UrlTarget:       report.Urls,
		PoC:             report.Description,
		MoreInfo:        "", // TODO: Extract from tags
		JudgeUser:       "",
		Amount:          0,
		CVSSJudge:       "",
		ScoreJudge:      "",
		DateTo:          report.DateTo,
		JudgeInfo:       "",
		LinkMoreInfo:    "",
		Attachment:      "",
		RavroVer:        "v2.0.0",
		Valuation:       "",
	}
}
