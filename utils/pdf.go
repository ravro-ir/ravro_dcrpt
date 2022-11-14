package utils

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"io/ioutil"
	"log"
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
			log.Fatal(errDir)
		}
	}
	err1 := ioutil.WriteFile(tmpPath+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err1 != nil {
		panic(err1)
	}
	f, err := os.Open(tmpPath + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(f)
	}
	if err != nil {
		log.Fatal(err)
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}
	}(dir + tmpPath)

	return true, nil
}
