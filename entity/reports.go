package entity

type InfoReport []struct {
	InfoDescription string `json:"infoDescription"`
	InfoTitle       string `json:"infoTitle"`
	InfoSolution    string `json:"infoSolution"`
	MoreInfo        string `json:"infoMore"`
}

type Report struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Reproduce       string `json:"reproduce"`
	DateFrom        string `json:"dateFrom"`
	CVSS            string `json:"cvss"`
	HunterUsername  string `json:"hunterUsername"`
	CompanyUsername string `json:"companyUsername"`
	Slug            string `json:"slug"`
	SubmissionDate  string `json:"submissionDate"`
	Ips             string `json:"ips"`
	Attachment      bool
	Scenario        string `json:"scenario"`
	ReportInfo      InfoReport
}
