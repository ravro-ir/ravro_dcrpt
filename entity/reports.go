package entity

//type InfoReport []struct {
//	InfoDescription string `json:"infoDescription"`
//	InfoTitle       string `json:"infoTitle"`
//	InfoSolution    string `json:"infoSolution"`
//	MoreInfo        string `json:"infoMore"`
//}

type InfoReport struct {
	Tags []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Solution    string `json:"solution"`
		MoreInfo    string `json:"moreInfo"`
	} `json:"tags"`
	Details struct {
		Judges []struct {
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		} `json:"judges"`
		Attachments []struct {
			Filename string `json:"filename"`
			Type     string `json:"type"`
		} `json:"attachments"`
		Cvss struct {
			Hunter struct {
				Vector string      `json:"vector"`
				Score  interface{} `json:"score"`
			} `json:"hunter"`
			Judge struct {
				Vector interface{} `json:"vector"`
				Score  interface{} `json:"score"`
			} `json:"judge"`
			Final struct {
				Vector interface{} `json:"vector"`
				Score  interface{} `json:"score"`
			} `json:"final"`
		} `json:"cvss"`
	} `json:"details"`
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
