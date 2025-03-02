//go:build darwin
// +build darwin

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	wk "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// RequestPdf structure to manage PDF creation
type RequestPdf struct {
	body string
}

// NewRequestPdf creates a new PDF request
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

// ParseTemplate parses an HTML template with dynamic data
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()

	// Properly replace HTML entities
	r.body = strings.ReplaceAll(r.body, "&#34;&lt;", "<")
	r.body = strings.ReplaceAll(r.body, "&gt;", ">")
	r.body = strings.ReplaceAll(r.body, "&lt;", "<")
	r.body = strings.ReplaceAll(r.body, "&gt;", ">")
	r.body = strings.ReplaceAll(r.body, "&#34;", "")

	return nil
}

// GeneratePDF generates the PDF file from the parsed HTML using go-wkhtmltopdf
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {
	// Validate HTML content
	if len(strings.TrimSpace(r.body)) == 0 {
		return false, errors.New("HTML content is empty")
	}

	// Create temporary directory for HTML files if it does not exist (optional, for debugging)
	tmpPath := "template/"
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		if err := os.Mkdir(tmpPath, 0777); err != nil {
			return false, err
		}
	}

	// Optionally, write the HTML to a temp file for debugging purposes
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	htmlFileName := tmpPath + timestamp + ".html"
	if err := ioutil.WriteFile(htmlFileName, []byte(r.body), 0644); err != nil {
		return false, err
	}
	// Clean up the temporary file when done
	defer os.Remove(htmlFileName)

	// Check OS before creating the PDF generator
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return false, errors.New("unsupported operating system")
	}

	// Create a new PDF generator instance
	pdfg, err := wk.NewPDFGenerator()
	if err != nil {
		return false, fmt.Errorf("failed to create PDF generator: %w", err)
	}

	// Set global options (similar to command-line options)
	pdfg.Dpi.Set(300)
	pdfg.MarginBottom.Set(10)
	pdfg.MarginTop.Set(10)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)

	// Create a new page from the HTML string (you can also use NewPageReader with a reader)
	page := wk.NewPageReader(strings.NewReader(r.body))

	// Add the page to the PDF generator
	pdfg.AddPage(page)

	// Create the PDF document in memory
	if err = pdfg.Create(); err != nil {
		return false, fmt.Errorf("failed to create PDF: %w", err)
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(pdfPath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0777); err != nil {
			return false, err
		}
	}

	// Write the PDF to file
	if err = pdfg.WriteFile(pdfPath); err != nil {
		return false, fmt.Errorf("failed to write PDF file: %w", err)
	}
	return true, nil
}
