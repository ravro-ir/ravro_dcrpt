package models

import "time"

// Report represents a bug bounty report
type Report struct {
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Reproduce       string     `json:"reproduce"`
	Scenario        string     `json:"scenario"`
	HunterUsername  string     `json:"hunterUsername"`
	CompanyUsername string     `json:"companyUsername"`
	Urls            string     `json:"urls"`
	Ips             string     `json:"ips"`
	SubmissionDate  string     `json:"submissionDate"`
	DateFrom        string     `json:"dateFrom"`
	DateTo          string     `json:"dateTo"`
	ReportInfo      InfoReport `json:"report_info"`
}

// InfoReport contains additional report information
type InfoReport struct {
	Details ReportDetails `json:"details"`
	Tags    []TagInfo     `json:"tags"`
}

// ReportDetails contains detailed information about the report
type ReportDetails struct {
	Target        string       `json:"target"`
	CurrentStatus string       `json:"currentStatus"`
	Cvss          CVSSInfo     `json:"cvss"`
	Attachments   []Attachment `json:"attachments"`
	Judges        []Judge      `json:"judges"`
}

// CVSSInfo contains CVSS scoring information
type CVSSInfo struct {
	Hunter CVSSScore `json:"hunter"`
	Judge  CVSSScore `json:"judge"`
}

// CVSSScore represents a single CVSS score
type CVSSScore struct {
	Score  string `json:"score"`
	Vector string `json:"vector"`
	Rating string `json:"rating"`
}

// Attachment represents a file attachment
type Attachment struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

// Judge represents a judge/reviewer
type Judge struct {
	Name string `json:"name"`
}

// TagInfo contains vulnerability tag information
type TagInfo struct {
	InfoTitle       string `json:"info_title"`
	InfoDescription string `json:"info_description"`
	InfoSolution    string `json:"info_solution"`
	InfoMore        string `json:"info_more"`
}

// Judgment represents the judgment data for a report
type Judgment struct {
	Reward        int               `json:"reward"`
	Description   string            `json:"description"`
	Cvss          JudgmentCVSS      `json:"cvss"`
	Vulnerability VulnerabilityInfo `json:"vulnerability"`
}

// JudgmentCVSS contains judgment CVSS information
type JudgmentCVSS struct {
	Value  string `json:"value"`
	Rating string `json:"rating"`
}

// VulnerabilityInfo contains vulnerability details
type VulnerabilityInfo struct {
	Name    string `json:"name"`
	Define  string `json:"define"`
	Fix     string `json:"fix"`
	Writeup string `json:"writeup"`
}

// Amendment represents additional amendment information
type Amendment struct {
	ReportID         string `json:"reportId"`
	AttachmentsCount int    `json:"attachmentsCount"`
	Description      string `json:"description"`
	HunterUsername   string `json:"hunterUsername"`
	CompanyUsername  string `json:"companyUsername"`
	SubmissionDate   string `json:"submissionDate"`

	// Legacy fields for backward compatibility
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// PDF represents the complete data structure for PDF generation
type PDF struct {
	Report    Report      `json:"report"`
	Judge     Judgment    `json:"judge"`
	Amendment []Amendment `json:"amendment,omitempty"`
}
