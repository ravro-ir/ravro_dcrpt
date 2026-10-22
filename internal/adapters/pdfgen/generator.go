package pdfgen

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"ravro_dcrpt/internal/ports"
	"ravro_dcrpt/pkg/models"
)

// Generator writes a standalone HTML report using the frontend template.
type Generator struct{}

// NewHTMLGenerator creates an HTML report exporter.
func NewHTMLGenerator() ports.HTMLGenerator {
	return &Generator{}
}

// GenerateReport writes the frontend-styled HTML report to outputPath.
func (g *Generator) GenerateReport(report *models.Report, judgment *models.Judgment, amendments []models.Amendment, outputPath string) error {
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	if err := writeReportFonts(outputDir); err != nil {
		return fmt.Errorf("failed to write fonts: %w", err)
	}

	body := buildReportHTML(report, judgment, amendments, outputDir)
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("HTML content is empty")
	}

	body = injectFontFace(body)
	if err := os.WriteFile(outputPath, []byte(body), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}
	return nil
}

func esc(value string) string {
	return html.EscapeString(value)
}
