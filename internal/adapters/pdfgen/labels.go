package pdfgen

import "strings"

const (
	labelTopbarTitle          = "گزارش آسیب‌پذیری"
	labelPrint                = "چاپ / PDF"
	labelThemeLight           = "حالت روشن"
	labelThemeDark            = "حالت تیره"
	labelReportID             = "شناسه گزارش"
	labelHunter               = "شکارچی"
	labelSubmittedAt          = "تاریخ ثبت"
	labelActivityRange        = "بازه زمانی فعالیت"
	labelTarget               = "هدف (Target)"
	labelIPs                  = "IP شکارچی"
	labelURLs                 = "آدرس‌های هدف (URLs)"
	labelCVSSSection          = "امتیازدهی CVSS"
	labelCVSSRoleCol          = "نقش"
	labelCVSSVersionCol       = "نسخه"
	labelCVSSScoreCol         = "امتیاز"
	labelCVSSSeverityCol      = "درجه"
	labelCVSSVectorCol        = "وکتور"
	labelRoleHunter           = "شکارچی"
	labelRoleJudge            = "داور"
	labelRoleCompany          = "میدان"
	labelRoleFinal            = "نهایی"
	labelCVSSNullScore        = "در انتظار"
	labelCVSSNullVector       = "هنوز ارزیابی نشده"
	labelCVSSFinalNull        = "پس از بررسی داور و میدان مشخص می‌شود"
	labelScenario             = "سناریوی آسیب‌پذیری"
	labelSolution             = " راه حل"
	labelCompany              = "میدان"
	labelUpdatedAt            = "آخرین بروزرسانی"
	labelBountyJudge          = "پاداش داور"
	labelBountyCompany        = "پاداش میدان"
	labelBountyFinal          = "پاداش نهایی"
	labelBountyPaidAt         = "تاریخ پرداخت به شکارچی"
	labelBountyCurrency       = "ریال"
	labelDuplicateOf          = "گزارش تکراری"
	labelTargetVersion        = "نسخهٔ هدف"
	labelRulesTitle           = "قوانین میدان"
	labelTimelineTitle        = "خط زمانی گزارش"
	labelAppealTitle          = "تجدیدنظر"
	labelAppealType           = "نوع"
	labelAppealTypeCompany    = "اعتراض به میدان"
	labelAppealTypeJudge      = "اعتراض به داور"
	labelAppealMessage        = "متن شکارچی"
	labelAppealReview         = "پاسخ"
	labelAppealReviewedAt     = "تاریخ تصمیم"
	labelAppealPending        = "در انتظار"
	labelAppealAccepted       = "پذیرفته شده"
	labelAppealRejected       = "رد شده"
	labelAmendmentStage       = "تغییر وضعیت"
	labelAmendmentFrom        = "از"
	labelAmendmentTo          = "به"
	labelAmendmentAttachments = "فایل‌های پیوست"
	labelAttachmentsTitle     = "فایل‌های پیوست"
	labelAttachmentsNote      = "فایل‌های پیوست در کنار این فایل HTML در پوشه ZIP قرار دارند."
	labelAttachmentsCount     = "%count% فایل"
	labelCommentsTitle        = "گفتگو و پیگیری"
	labelCommentsCount        = "%count% پیام"
	labelFooterBrand          = "تولیدشده توسط پلتفرم راورو"
	labelAmendmentsTitle      = "اصلاحات"
	labelAmendmentsCount      = "%count% اصلاح"
	labelAmendmentNumber      = "اصلاح"
	labelAmendmentExplanation = "توضیحات"
	labelAmendmentDate        = "تاریخ"
)

func severityLabelFA(level string) string {
	switch strings.ToUpper(level) {
	case "CRITICAL":
		return "بحرانی"
	case "HIGH":
		return "بالا"
	case "MEDIUM":
		return "متوسط"
	case "LOW":
		return "پایین"
	case "NONE":
		return "ندارد"
	default:
		return "نامشخص"
	}
}

func commentRoleLabelFA(role string) string {
	switch strings.ToUpper(role) {
	case "COMPANY":
		return "میدان"
	case "JUDGE":
		return "داور"
	case "SUPPORT":
		return "پشتیبانی راورو"
	default:
		return "شکارچی"
	}
}

func statusLabelFA(status string) string {
	switch strings.ToUpper(status) {
	case "COMPANY":
		return "در بررسی میدان"
	case "JUDGE":
		return "در بررسی داور"
	case "RESOLVED":
		return "حل‌شده"
	case "REJECTED":
		return "رد شده"
	case "PENDING":
		return "در انتظار بررسی"
	case "INITIAL":
		return "بررسی اولیه"
	case "INITIAL_REJECT":
		return "نیازمند اطلاعات بیشتر"
	case "INITIAL_REVIEW":
		return "بازبینی اولیه"
	case "INITIAL_REJECT_SPAM":
		return "رد شده (اسپم)"
	case "JUDGEMENT":
		return "بررسی تیم داوری"
	case "JUDGEMENT_REJECT":
		return "نیازمند اطلاعات بیشتر"
	case "JUDGEMENT_REVIEW":
		return "بازبینی مجدد تیم داوری"
	case "JUDGEMENT_REJECT_CLOSE":
		return "رد شده (تیم داوری)"
	case "JUDGEMENT_REJECT_INCOMPLETE":
		return "ناقص"
	case "RAVRO":
		return "ارزیابی نهایی"
	case "RAVRO_REVISION_REQUEST":
		return "درخواست تجدیدنظر"
	case "RAVRO_REVISION_REVIEW":
		return "بررسی درخواست تجدیدنظر"
	case "RAVRO_OPINION":
		return "اعلام نظر راورو"
	case "ACCEPTED":
		return "تایید گزارش"
	case "DUPLICATE":
		return "تکراری"
	case "INFORMATIVE":
		return "آموزنده"
	case "REJECTED_SPAM":
		return "اسپم"
	case "REJECTED_INCOMPLETE":
		return "ناقص"
	case "REJECTED_OUT_OF_SCOPE":
		return "خارج از محدوده"
	case "RETEST", "RETEST_REQUEST":
		return "درخواست ریتست"
	case "DISCLOSE":
		return "درخواست افشای آسیب پذیری"
	case "DISCLOSE_REQUEST":
		return "بررسی درخواست افشای آسیب پذیری"
	case "DISCLOSED":
		return "افشا شده (درخواست افشای آسیب پذیری)"
	default:
		return status
	}
}
