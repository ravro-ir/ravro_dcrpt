package ports

import "ravro_dcrpt/pkg/models"

// PDFGenerator defines the interface for PDF generation
type PDFGenerator interface {
	// GenerateReport creates a PDF from report data
	GenerateReport(report *models.Report, judgment *models.Judgment, outputPath string) error

	// GenerateFromHTML creates a PDF from HTML template
	GenerateFromHTML(htmlContent string, outputPath string) error

	// SetTemplate sets the HTML template for report generation
	SetTemplate(templatePath string) error
}
