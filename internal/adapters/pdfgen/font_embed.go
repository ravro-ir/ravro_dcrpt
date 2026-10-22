package pdfgen

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
)

//go:embed fonts/Vazirmatn-Regular.ttf
var vazirmatnRegular []byte

//go:embed fonts/Vazirmatn-Bold.ttf
var vazirmatnBold []byte

const vazirmatnFaceCSS = `
        @font-face {
            font-family: 'Vazirmatn';
            src: url('Vazirmatn-Regular.ttf') format('truetype');
            font-weight: 300;
            font-style: normal;
        }
        @font-face {
            font-family: 'Vazirmatn';
            src: url('Vazirmatn-Regular.ttf') format('truetype');
            font-weight: 400;
            font-style: normal;
        }
        @font-face {
            font-family: 'Vazirmatn';
            src: url('Vazirmatn-Bold.ttf') format('truetype');
            font-weight: 500;
            font-style: normal;
        }
        @font-face {
            font-family: 'Vazirmatn';
            src: url('Vazirmatn-Bold.ttf') format('truetype');
            font-weight: 600;
            font-style: normal;
        }
        @font-face {
            font-family: 'Vazirmatn';
            src: url('Vazirmatn-Bold.ttf') format('truetype');
            font-weight: 700;
            font-style: normal;
        }
        @font-face {
            font-family: 'Vazirmatn';
            src: url('Vazirmatn-Bold.ttf') format('truetype');
            font-weight: bold;
            font-style: normal;
        }
`

func writeReportFonts(dir string) error {
	if err := os.WriteFile(filepath.Join(dir, "Vazirmatn-Regular.ttf"), vazirmatnRegular, 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "Vazirmatn-Bold.ttf"), vazirmatnBold, 0644)
}

func injectFontFace(html string) string {
	return strings.Replace(html, "<style>", "<style>"+vazirmatnFaceCSS, 1)
}
