package pdfgen

import "testing"

func TestStatusLabelFAMatchesFrontendHTML(t *testing.T) {
	// htmlReport.statusLabels overlay STAGES from bugbounty-reports.json (FA).
	want := map[string]string{
		"COMPANY":                     "در بررسی میدان",
		"JUDGE":                       "در بررسی داور",
		"RESOLVED":                    "حل‌شده",
		"REJECTED":                    "رد شده",
		"PENDING":                     "در انتظار بررسی",
		"INITIAL":                     "بررسی اولیه",
		"INITIAL_REJECT":              "نیازمند اطلاعات بیشتر",
		"INITIAL_REVIEW":              "بازبینی اولیه",
		"INITIAL_REJECT_SPAM":         "رد شده (اسپم)",
		"JUDGEMENT":                   "بررسی تیم داوری",
		"JUDGEMENT_REJECT":            "نیازمند اطلاعات بیشتر",
		"JUDGEMENT_REVIEW":            "بازبینی مجدد تیم داوری",
		"JUDGEMENT_REJECT_CLOSE":      "رد شده (تیم داوری)",
		"JUDGEMENT_REJECT_INCOMPLETE": "ناقص",
		"RAVRO":                       "ارزیابی نهایی",
		"RAVRO_REVISION_REQUEST":      "درخواست تجدیدنظر",
		"RAVRO_REVISION_REVIEW":       "بررسی درخواست تجدیدنظر",
		"RAVRO_OPINION":               "اعلام نظر راورو",
		"ACCEPTED":                    "تایید گزارش",
		"DUPLICATE":                   "تکراری",
		"INFORMATIVE":                 "آموزنده",
		"REJECTED_SPAM":               "اسپم",
		"REJECTED_INCOMPLETE":         "ناقص",
		"REJECTED_OUT_OF_SCOPE":       "خارج از محدوده",
		"RETEST":                      "درخواست ریتست",
		"RETEST_REQUEST":              "درخواست ریتست",
		"DISCLOSE":                    "درخواست افشای آسیب پذیری",
		"DISCLOSE_REQUEST":            "بررسی درخواست افشای آسیب پذیری",
		"DISCLOSED":                   "افشا شده (درخواست افشای آسیب پذیری)",
	}
	for status, label := range want {
		if got := statusLabelFA(status); got != label {
			t.Errorf("statusLabelFA(%q)=%q, want %q", status, got, label)
		}
	}
}

func TestStaticCopyMatchesFrontendFA(t *testing.T) {
	if labelSolution != " راه حل" {
		t.Fatalf("solution title: got %q", labelSolution)
	}
	if labelAttachmentsNote != "فایل‌های پیوست در کنار این فایل HTML در پوشه ZIP قرار دارند." {
		t.Fatalf("attachments note: got %q", labelAttachmentsNote)
	}
	if labelScenario != "سناریوی آسیب‌پذیری" {
		t.Fatalf("scenario title: got %q", labelScenario)
	}
}
