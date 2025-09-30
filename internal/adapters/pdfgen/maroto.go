package pdfgen

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"ravro_dcrpt/internal/ports"
	"ravro_dcrpt/pkg/models"
)

// MarotoGenerator implements PDFGenerator using Maroto library
type MarotoGenerator struct {
	templatePath string
}

// NewMarotoGenerator creates a new Maroto-based PDF generator
func NewMarotoGenerator() ports.PDFGenerator {
	return &MarotoGenerator{}
}

// SetTemplate sets the HTML template (not used in Maroto, kept for interface compatibility)
func (g *MarotoGenerator) SetTemplate(templatePath string) error {
	g.templatePath = templatePath
	return nil
}

// GenerateReport creates a PDF from report data
func (g *MarotoGenerator) GenerateReport(report *models.Report, judgment *models.Judgment, outputPath string) error {
	// Configure PDF with RTL support for Persian text
	cfg := config.NewBuilder().
		WithMargins(10, 10, 10).
		Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	// Add content to PDF
	err := g.buildPDFContent(m, report, judgment)
	if err != nil {
		return fmt.Errorf("failed to build PDF content: %w", err)
	}

	// Generate PDF file
	document, err := m.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save to file
	err = document.Save(outputPath)
	if err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}

	return nil
}

// buildPDFContent builds the PDF content structure
func (g *MarotoGenerator) buildPDFContent(m core.Maroto, report *models.Report, judgment *models.Judgment) error {
	// Add header
	g.addHeader(m, report)

	// Add report details
	g.addReportDetails(m, report)

	// Add vulnerability information
	if judgment != nil {
		g.addJudgmentInfo(m, judgment)
	}

	// Add description/PoC
	g.addDescription(m, report)

	// Add attachments
	g.addAttachments(m, report)

	return nil
}

// addHeader adds the PDF header
func (g *MarotoGenerator) addHeader(m core.Maroto, report *models.Report) {
	m.AddPages(
		page.New().Add(
			row.New(20).Add(
				col.New(12).Add(
					text.New("گزارش آسیب‌پذیری - Ravro Platform", props.Text{
						Top:   5,
						Size:  16,
						Align: align.Center,
						Style: fontstyle.Bold,
					}),
				),
			),
			row.New(10).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("شناسه گزارش: %s", report.Slug), props.Text{
						Top:   2,
						Size:  12,
						Align: align.Center,
					}),
				),
			),
		),
	)
}

// addReportDetails adds report details section
func (g *MarotoGenerator) addReportDetails(m core.Maroto, report *models.Report) {
	m.AddPages(
		page.New().Add(
			row.New(10).Add(
				col.New(12).Add(
					text.New("عنوان: "+report.Title, props.Text{
						Size:  12,
						Style: fontstyle.Bold,
						Top:   2,
					}),
				),
			),
			row.New(8).Add(
				col.New(6).Add(
					text.New("شکارچی: "+report.HunterUsername, props.Text{
						Size: 10,
						Top:  2,
					}),
				),
				col.New(6).Add(
					text.New("شرکت: "+report.CompanyUsername, props.Text{
						Size: 10,
						Top:  2,
					}),
				),
			),
			row.New(8).Add(
				col.New(6).Add(
					text.New("تاریخ ثبت: "+report.SubmissionDate, props.Text{
						Size: 10,
						Top:  2,
					}),
				),
				col.New(6).Add(
					text.New("وضعیت: "+report.ReportInfo.Details.CurrentStatus, props.Text{
						Size: 10,
						Top:  2,
					}),
				),
			),
		),
	)
}

// addJudgmentInfo adds judgment information
func (g *MarotoGenerator) addJudgmentInfo(m core.Maroto, judgment *models.Judgment) {
	m.AddPages(
		page.New().Add(
			row.New(10).Add(
				col.New(12).Add(
					text.New("اطلاعات داوری", props.Text{
						Size:  14,
						Style: fontstyle.Bold,
						Top:   3,
					}),
				),
			),
			row.New(8).Add(
				col.New(6).Add(
					text.New(fmt.Sprintf("پاداش: %d تومان", judgment.Reward), props.Text{
						Size: 10,
						Top:  2,
					}),
				),
				col.New(6).Add(
					text.New("امتیاز CVSS: "+judgment.Cvss.Rating, props.Text{
						Size: 10,
						Top:  2,
					}),
				),
			),
			row.New(8).Add(
				col.New(12).Add(
					text.New("نوع آسیب‌پذیری: "+judgment.Vulnerability.Name, props.Text{
						Size: 10,
						Top:  2,
					}),
				),
			),
		),
	)
}

// addDescription adds description/PoC section
func (g *MarotoGenerator) addDescription(m core.Maroto, report *models.Report) {
	m.AddPages(
		page.New().Add(
			row.New(10).Add(
				col.New(12).Add(
					text.New("توضیحات و PoC", props.Text{
						Size:  14,
						Style: fontstyle.Bold,
						Top:   3,
					}),
				),
			),
			row.New(30).Add(
				col.New(12).Add(
					text.New(g.stripHTML(report.Description), props.Text{
						Size: 9,
						Top:  2,
					}),
				),
			),
		),
	)
}

// addAttachments adds attachments section
func (g *MarotoGenerator) addAttachments(m core.Maroto, report *models.Report) {
	if len(report.ReportInfo.Details.Attachments) == 0 {
		return
	}

	m.AddPages(
		page.New().Add(
			row.New(10).Add(
				col.New(12).Add(
					text.New("فایل‌های پیوست", props.Text{
						Size:  14,
						Style: fontstyle.Bold,
						Top:   3,
					}),
				),
			),
		),
	)

	for _, att := range report.ReportInfo.Details.Attachments {
		m.AddPages(
			page.New().Add(
				row.New(6).Add(
					col.New(12).Add(
						text.New("• "+att.Filename, props.Text{
							Size: 9,
							Top:  1,
						}),
					),
				),
			),
		)
	}
}

// stripHTML removes HTML tags from text (basic implementation)
func (g *MarotoGenerator) stripHTML(html string) string {
	// Basic HTML stripping - can be improved with proper HTML parser
	result := strings.ReplaceAll(html, "<br>", "\n")
	result = strings.ReplaceAll(result, "<br/>", "\n")
	result = strings.ReplaceAll(result, "<br />", "\n")
	result = strings.ReplaceAll(result, "</p>", "\n")

	// Remove all other HTML tags
	inTag := false
	var sb strings.Builder
	for _, r := range result {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// GenerateFromHTML creates a PDF from HTML (simplified for Pure Go)
func (g *MarotoGenerator) GenerateFromHTML(htmlContent string, outputPath string) error {
	// For now, we'll use a simple text conversion
	// In a full implementation, you might want to use a proper HTML parser

	cfg := config.NewBuilder().
		WithMargins(10, 10, 10).
		Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	// Convert HTML to plain text and add to PDF
	plainText := g.stripHTML(htmlContent)

	m.AddPages(
		page.New().Add(
			row.New(200).Add(
				col.New(12).Add(
					text.New(plainText, props.Text{
						Size: 10,
						Top:  5,
					}),
				),
			),
		),
	)

	document, err := m.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate PDF from HTML: %w", err)
	}

	return document.Save(outputPath)
}

// GenerateHTMLReport generates an HTML file from report data
func (g *MarotoGenerator) GenerateHTMLReport(report *models.Report, judgment *models.Judgment, outputPath string, templateData interface{}) error {
	// Parse the HTML template
	tmpl, err := template.New("report").Parse(HTMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Replace .pdf extension with .html
	htmlOutputPath := strings.Replace(outputPath, ".pdf", ".html", 1)

	// Create output file
	f, err := os.Create(htmlOutputPath)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer f.Close()

	// Execute template
	if err := tmpl.Execute(f, templateData); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return nil
}
