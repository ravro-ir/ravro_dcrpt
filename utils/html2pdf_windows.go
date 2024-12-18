//go:build windows
// +build windows

package utils

import (
	"bytes"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// RequestPdf pdf requestpdf struct
type RequestPdf struct {
	body string
}

// NewRequestPdf new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

// ParseTemplate parsing template function
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

func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {
	// Initialize library
	if err := pdf.Init(); err != nil {
		return false, err
	}
	defer pdf.Destroy()
	tmpPath := "template/"
	t := time.Now().Unix()

	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		if err := os.Mkdir(tmpPath, 0777); err != nil {
			return false, err
		}
	}
	htmlFileName := tmpPath + strconv.FormatInt(int64(t), 10) + ".html"
	if err := ioutil.WriteFile(htmlFileName, []byte(r.body), 0644); err != nil {
		return false, err
	}
	object, err := pdf.NewObject(htmlFileName)
	if err != nil {
		return false, err
	}
	converter, err := pdf.NewConverter()
	if err != nil {
		return false, err
	}
	defer converter.Destroy()
	converter.Add(object)
	outFile, err := os.Create(pdfPath)
	if err != nil {
		return false, err
	}
	defer outFile.Close()
	if err := converter.Run(outFile); err != nil {
		return false, err
	}
	if err := os.Remove(htmlFileName); err != nil {
		return false, err
	}

	return true, nil
}
