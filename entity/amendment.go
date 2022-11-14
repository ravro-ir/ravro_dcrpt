package entity

type Amendment struct {
	AttachmentsCount int    `json:"attachmentsCount"`
	Description      string `json:"description"`
	HunterUsername   string `json:"hunterUsername"`
	CompanyUsername  string `json:"companyUsername"`
	SubmissionDate   string `json:"submissionDate"`
}
