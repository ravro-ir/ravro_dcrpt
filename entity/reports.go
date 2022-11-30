package entity

type InfoReport struct {
	Tags []struct {
		InfoTitle       string `json:"infoTitle"`
		InfoDescription string `json:"infoDescription"`
		InfoSolution    string `json:"infoSolution"`
		InfoMore        string `json:"infoMore"`
	} `json:"tags"`
	Details struct {
		CurrentStatus string `json:"currentStatus"`
		Target        string `json:"target"`
		Judges        []struct {
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		} `json:"judges"`
		Attachments []struct {
			Filename string `json:"filename"`
			Type     string `json:"type"`
		} `json:"attachments"`
		Cvss struct {
			Hunter struct {
				Vector string `json:"vector"`
				Score  string `json:"score"`
			} `json:"hunter"`
			Judge struct {
				Vector string `json:"vector"`
				Score  string `json:"score"`
			} `json:"judge"`
			Company struct {
				Vector string `json:"vector"`
				Score  string `json:"score"`
			} `json:"company"`
			Final struct {
				Vector string `json:"vector"`
				Score  string `json:"score"`
			} `json:"final"`
		} `json:"cvss"`
	} `json:"details"`
}

type Report struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	Reproduce        string `json:"reproduce"`
	Scenario         string `json:"scenario"`
	SubmissionDate   string `json:"submissionDate"`
	HunterUsername   string `json:"hunterUsername"`
	CompanyUsername  string `json:"companyUsername"`
	AttachmentsCount int    `json:"attachmentsCount"`
	Slug             string `json:"slug"`
	Category         struct {
		Title  string `json:"title"`
		Impact string `json:"impact"`
	} `json:"category"`
	Ips        string `json:"ips"`
	Cvss       string `json:"cvss"`
	Urls       string `json:"urls"`
	DateFrom   string `json:"dateFrom"`
	DateTo     string `json:"dateTo"`
	ReportInfo InfoReport
}
