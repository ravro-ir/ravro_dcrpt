package pdfgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ravro_dcrpt/pkg/models"
)

func sampleReport() *models.Report {
	return &models.Report{
		Title:           "دسترسی عمومی به فایل‌های حساس با انجام تنظیمات اشتباه",
		Slug:            "r2026-09-06-0002",
		HunterUsername:  "mahyarm02",
		CompanyUsername: "basalam",
		Description:     "برای رفع مشکل باید دسترسی عمومی را محدود کرد.\n\n```nginx\nlocation ~ /\\.(?!well-known) {\n    deny all;\n}\n```",
		Scenario:        "در حین بررسی مشاهده شد که فایل‌های حساس بدون احراز هویت در دسترس هستند.\n\nhttps://experience.basalam.com/.DS_Store",
		Urls:            "https://experience.basalam.com/.DS_Store,https://experience.basalam.com/web.config",
		Ips:             "141.11.1.208",
		Cvss:            "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:N/A:N",
		DateFrom:        "2026-09-05",
		DateTo:          "2026-09-06",
		SubmissionDate:  "2026-09-06 10:00:00",
		TargetName:      "وب سایت باسلام",
		ReportInfo: models.InfoReport{
			Details: models.ReportDetails{
				Target:        "وب سایت باسلام",
				CurrentStatus: "JUDGEMENT",
				Cvss: models.CVSSInfo{
					Hunter: models.CVSSScore{Vector: "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:N/A:N", Score: "5.3"},
				},
				Attachments: []models.Attachment{
					{Filename: "files.zip"},
					{Filename: "Screenshot_217.png"},
				},
			},
		},
	}
}

func TestGenerateReportWritesFrontendHTML(t *testing.T) {
	out := filepath.Join(t.TempDir(), "basalam__r2026-09-06-0002__mahyarm02.html")
	gen := NewHTMLGenerator()
	report := sampleReport()

	if err := gen.GenerateReport(report, nil, nil, out); err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(out) != ".html" {
		t.Fatalf("expected .html output, got %s", out)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() < 2000 {
		t.Fatalf("html too small: %d", info.Size())
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(out), "Vazirmatn-Regular.ttf")); err != nil {
		t.Fatalf("expected embedded font next to HTML: %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	html := string(raw)

	for _, needle := range []string{
		`lang="fa"`,
		`dir="rtl"`,
		`class="hero"`,
		`class="info-grid"`,
		`class="cvss-table"`,
		`class="attachments-grid"`,
		`class="status-pill status-judge`,
		`class="severity-block sev-medium"`,
		"سناریوی آسیب‌پذیری",
		"امتیازدهی CVSS",
		"دسترسی عمومی به فایل‌های حساس",
		"وب سایت باسلام",
		"mahyarm02",
		"r2026-09-06-0002",
		"141.11.1.208",
		"https://experience.basalam.com/.DS_Store",
		"location ~",
		"files.zip",
		"Screenshot_217.png",
		"@font-face",
		"Vazirmatn-Regular.ttf",
		"theme-toggle",
		" راه حل",
		"فایل‌های پیوست در کنار این فایل HTML در پوشه ZIP قرار دارند.",
		".md-content strong, .md-content b { font-weight: 700; }",
		"font-weight: 500;",
		"font-weight: bold;",
	} {
		if !strings.Contains(html, needle) {
			t.Fatalf("expected HTML to contain %q", needle)
		}
	}
	if strings.Contains(html, "fonts.googleapis.com") {
		t.Fatal("HTML should not load remote Google Fonts")
	}
}

func TestGenerateReportEscapesUntrustedFields(t *testing.T) {
	out := filepath.Join(t.TempDir(), "co__id__h.html")
	report := sampleReport()
	report.Title = `<script>alert(1)</script>`
	report.HunterUsername = `hunter"><img src=x>`
	report.Urls = `javascript:alert(1)`

	if err := NewHTMLGenerator().GenerateReport(report, nil, nil, out); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, out)
	if strings.Contains(html, "<script>alert(1)</script>") {
		t.Fatal("title was not escaped")
	}
	if !strings.Contains(html, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatal("expected escaped title")
	}
	if strings.Contains(html, `hunter"><img`) {
		t.Fatal("hunter was not escaped")
	}
}

func TestGenerateReportConditionalCommentsAndAmendments(t *testing.T) {
	out := filepath.Join(t.TempDir(), "co__id__h.html")
	report := sampleReport()

	html := readGenerated(t, report, nil, nil, out)
	if strings.Contains(html, "گفتگو و پیگیری") {
		t.Fatal("comments section should be omitted when empty")
	}
	if strings.Contains(html, "اصلاحات") {
		t.Fatal("amendments section should be omitted when empty")
	}

	report.Comments = []models.Comment{{
		Message:   "لطفا **جزئیات** بیشتری بفرستید",
		By:        models.CommentAuthor{Nickname: "judge1", Role: "JUDGE"},
		CreatedAt: "2026-09-07T10:00:00Z",
	}}
	amendments := []models.Amendment{{
		Description:    "توضیح اصلاح با `code`",
		CreatedAt:      time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
		SubmissionDate: "2026-09-08",
	}}
	html = readGenerated(t, report, nil, amendments, out)
	if !strings.Contains(html, "گفتگو و پیگیری") {
		t.Fatal("expected comments section")
	}
	if !strings.Contains(html, "اصلاحات") {
		t.Fatal("expected amendments section")
	}
	if !strings.Contains(html, "comment-role-judge") {
		t.Fatal("expected judge comment role class")
	}
	if !strings.Contains(html, "judge1") {
		t.Fatal("expected comment author nickname")
	}
	if !strings.Contains(html, "<strong>جزئیات</strong>") && !strings.Contains(html, "<strong>جزئیات</strong>") {
		if !strings.Contains(html, "جزئیات") {
			t.Fatal("expected markdown comment content")
		}
	}
	if !strings.Contains(html, "<code>code</code>") && !strings.Contains(html, "code") {
		t.Fatal("expected markdown in amendment")
	}

	report.Comments = []models.Comment{{
		Message:   "پیام پشتیبانی",
		By:        models.CommentAuthor{DisplayName: "wwrbvettony4", Role: "SUPPORT"},
		CreatedAt: "2026-10-01T08:00:00Z",
	}}
	html = readGenerated(t, report, nil, nil, out)
	if !strings.Contains(html, `class="comment-author">wwrbvettony4</span>`) {
		t.Fatal("expected display_name as comment author")
	}
	if strings.Contains(html, `class="comment-author">-</span>`) {
		t.Fatal("comment author should not fall back to dash")
	}

	report.Comments = []models.Comment{{
		Message:   "بدون نام",
		By:        models.CommentAuthor{Role: "SUPPORT"},
		CreatedAt: "2026-10-01T08:00:00Z",
	}}
	html = readGenerated(t, report, nil, nil, out)
	if !strings.Contains(html, `class="comment-author">پشتیبانی راورو</span>`) {
		t.Fatal("expected role label when comment author name is missing")
	}
}

func TestGenerateReportNumberedAttachmentHrefs(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "basalam__r2026-09-06-0002__mahyarm02.html")
	reportDir := filepath.Join(dir, "report")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reportDir, "001. files.zip"), []byte("zip"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reportDir, "002. Screenshot_217.png"), []byte("png"), 0644); err != nil {
		t.Fatal(err)
	}

	html := readGenerated(t, sampleReport(), nil, nil, out)
	if !strings.Contains(html, `href="report/001. files.zip"`) {
		t.Fatal("expected numbered zip href matching on-disk name")
	}
	if !strings.Contains(html, `href="report/002. Screenshot_217.png"`) {
		t.Fatal("expected numbered screenshot href")
	}
	if strings.Contains(html, `href="report/files.zip"`) {
		t.Fatal("should not link the unprefixed report_info filename")
	}
	if strings.Contains(html, "%20") {
		t.Fatal("local file hrefs must not percent-encode spaces")
	}
	if strings.Contains(html, `class="attachment-card" target="_blank"`) {
		t.Fatal("local attachment links must not open a new tab")
	}
}

func TestGenerateReportJudgmentFillsCVSSRow(t *testing.T) {
	out := filepath.Join(t.TempDir(), "co__id__h.html")
	report := sampleReport()
	judgment := &models.Judgment{
		Cvss: models.JudgmentCVSS{Value: "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:N", Rating: "8.1"},
	}
	html := readGenerated(t, report, judgment, nil, out)
	if !strings.Contains(html, "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:N") {
		t.Fatal("expected judge CVSS vector")
	}
	if !strings.Contains(html, "8.1") {
		t.Fatal("expected judge score")
	}
}

func TestGenerateReportRendersBackendExportMetadata(t *testing.T) {
	out := filepath.Join(t.TempDir(), "co__id__h.html")
	report := sampleReport()
	report.CompanyUsername = "acme"
	report.ReportInfo.Export = models.ReportInfoExport{
		Present:               true,
		UpdatedAt:             "2026-10-01T10:00:00Z",
		Bughunter:             models.ReportInfoParty{Title: "hunter-one", URL: "https://ravro.ir/fa/hunters/h1", Avatar: "hunter-avatar-code"},
		Company:               models.ReportInfoParty{Title: "Acme Field", URL: "https://ravro.ir/fa/companies/acme", Avatar: "company-avatar-code"},
		Target:                models.ReportInfoExportTarget{Title: "Shop", Version: "3", Rules: "دامنه باید داخلی باشد"},
		Bounty:                models.ReportInfoExportBounty{Judge: "1000000", Final: "1100000"},
		DuplicateOfReportSlug: "r2026-10-01-0004",
		Timeline: []models.ReportInfoTimelineEntry{{
			At: "2026-10-01T11:00:00Z", FromStage: "INITIAL", ToStage: "JUDGEMENT",
		}},
		Appeal: &models.ReportInfoAppeal{
			Status:     "PENDING",
			TargetRole: "COMPANY",
			Message:    "please review",
		},
	}
	amendments := []models.Amendment{{
		Description: "fix",
		FromStage:   "INITIAL",
		ToStage:     "JUDGEMENT",
		CreatedAt:   time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC),
	}}

	html := readGenerated(t, report, nil, amendments, out)
	for _, needle := range []string{
		"میدان",
		"Acme Field",
		"پاداش داور",
		"1,000,000 ریال",
		"گزارش تکراری",
		"r2026-10-01-0004",
		"نسخهٔ هدف",
		"قوانین میدان",
		"خط زمانی گزارش",
		"تجدیدنظر",
		"اعتراض به میدان",
		"please review",
		"تغییر وضعیت",
		"party-chip",
		"timeline-row",
		"appeal-grid",
		`https://api.ravro.ir/api/v1/public/attachments/avatar?code=hunter-avatar-code`,
		`https://api.ravro.ir/api/v1/public/attachments/avatar?code=company-avatar-code`,
		"تولیدشده توسط پلتفرم راورو",
	} {
		if !strings.Contains(html, needle) {
			t.Fatalf("expected HTML to contain %q", needle)
		}
	}
	if strings.Contains(html, "پلتفرم Ravro") {
		t.Fatal("footer should use راورو, not Ravro")
	}
	if strings.Contains(html, "fonts.googleapis.com") {
		t.Fatal("HTML should stay offline")
	}
}

func readGenerated(t *testing.T, report *models.Report, judgment *models.Judgment, amendments []models.Amendment, out string) string {
	t.Helper()
	if err := NewHTMLGenerator().GenerateReport(report, judgment, amendments, out); err != nil {
		t.Fatal(err)
	}
	return readFile(t, out)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}
