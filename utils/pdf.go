package utils

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}

//new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

//parsing template function
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
	tmpPath := "template/"
	t := time.Now().Unix()
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		errDir := os.Mkdir(tmpPath, 0777)
		if errDir != nil {
			return false, err
		}
	}
	err1 := ioutil.WriteFile(tmpPath+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err1 != nil {
		panic(err1)
	}
	f, err := os.Open(tmpPath + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer func(f *os.File) (bool, error) {
			err := f.Close()
			if err != nil {
				return false, err
			}
			return true, nil
		}(f)

	}
	if err != nil {
		return false, err
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return false, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return false, err
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		return false, err
	}

	dir, err := os.Getwd()
	if err != nil {
		return false, err
	}

	defer func(path string) (bool, error) {
		err := os.RemoveAll(path)
		if err != nil {
			return false, err
		}
		return true, nil
	}(dir + tmpPath)

	return true, nil
}
