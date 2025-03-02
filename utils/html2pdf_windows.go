//go:build windows
// +build windows

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	wk "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// RequestPdf manages PDF creation
type RequestPdf struct {
	body string
}

// NewRequestPdf creates a new instance
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{body: body}
}

// ParseTemplate loads and processes an HTML template
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
	r.body = strings.ReplaceAll(r.body, "&#34;&lt;", "<")
	r.body = strings.ReplaceAll(r.body, "&gt;", ">")
	r.body = strings.ReplaceAll(r.body, "&lt;", "<")
	r.body = strings.ReplaceAll(r.body, "&gt;", ">")
	r.body = strings.ReplaceAll(r.body, "&#34;", "")
	return nil
}

// GeneratePDF generates a PDF from the HTML
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {
	if len(strings.TrimSpace(r.body)) == 0 {
		return false, errors.New("HTML content is empty")
	}

	// Ensure wkhtmltopdf is accessible
	_, err := exec.LookPath("wkhtmltopdf.exe")
	if err != nil {
		return false, fmt.Errorf("wkhtmltopdf not found: %w", err)
	}

	tmpPath := "template/"
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		if err := os.Mkdir(tmpPath, 0777); err != nil {
			return false, err
		}
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	htmlFileName := tmpPath + timestamp + ".html"
	if err := ioutil.WriteFile(htmlFileName, []byte(r.body), 0644); err != nil {
		return false, err
	}
	defer os.Remove(htmlFileName)

	// Initialize wkhtmltopdf
	pdfg, err := wk.NewPDFGenerator()
	if err != nil {
		return false, fmt.Errorf("failed to create PDF generator: %w", err)
	}

	pdfg.Dpi.Set(300)
	pdfg.MarginBottom.Set(10)
	pdfg.MarginTop.Set(10)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)

	page := wk.NewPageReader(strings.NewReader(r.body))
	pdfg.AddPage(page)

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

	// Write the PDF to the specified file
	if err = pdfg.WriteFile(pdfPath); err != nil {
		return false, fmt.Errorf("failed to write PDF file: %w", err)
	}

	return true, nil
}
