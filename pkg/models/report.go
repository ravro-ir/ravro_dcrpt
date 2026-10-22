package models

import "time"

// Schema versions for encrypted report JSON.
const (
	SchemaVersionV1 = 1
	SchemaVersionV2 = 2
)

// Report represents a bug bounty report
type Report struct {
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Solution        string     `json:"solution,omitempty"`
	Reproduce       string     `json:"reproduce"`
	Scenario        string     `json:"scenario"`
	HunterUsername  string     `json:"hunterUsername"`
	CompanyUsername string     `json:"companyUsername"`
	Urls            string     `json:"urls"`
	Ips             string     `json:"ips"`
	SubmissionDate  string     `json:"submissionDate"`
	CreatedAt       string     `json:"createdAt,omitempty"`
	DateFrom        string     `json:"dateFrom"`
	DateTo          string     `json:"dateTo"`
	TargetName      string     `json:"targetName,omitempty"`
	Cvss            string     `json:"cvss,omitempty"`
	ReportInfo      InfoReport `json:"report_info"`
	SchemaVersion   int        `json:"-"`
	Comments        []Comment  `json:"-"`
}

// InfoReport contains additional report information
type InfoReport struct {
	SchemaVersion int              `json:"schema_version,omitempty"`
	Details       ReportDetails    `json:"details"`
	Tags          []TagInfo        `json:"tags"`
	Export        ReportInfoExport `json:"export"`
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
	Hunter  CVSSScore `json:"hunter"`
	Judge   CVSSScore `json:"judge"`
	Company CVSSScore `json:"company"`
	Final   CVSSScore `json:"final"`
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
	FromStage        string `json:"from_stage,omitempty"`
	ToStage          string `json:"to_stage,omitempty"`

	// Legacy fields for backward compatibility
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ReportInfoExport is schema v2 plaintext metadata from backend-api report_info.json.
type ReportInfoExport struct {
	Present               bool
	UpdatedAt             string
	Bughunter             ReportInfoParty
	Company               ReportInfoParty
	Target                ReportInfoExportTarget
	Bounty                ReportInfoExportBounty
	DuplicateOfReportSlug string
	Timeline              []ReportInfoTimelineEntry
	Appeal                *ReportInfoAppeal
}

// ReportInfoParty is hunter or company identity in the export block.
type ReportInfoParty struct {
	Title  string
	Slug   string
	URL    string
	Avatar string
}

// ReportInfoExportTarget is the target snapshot in the export block.
type ReportInfoExportTarget struct {
	Title   string
	Slug    string
	URL     string
	Version string
	Rules   string
}

// ReportInfoExportBounty holds proposal and payout amounts as numeric strings.
type ReportInfoExportBounty struct {
	Judge   string
	Company string
	Final   string
	PaidAt  string
}

// ReportInfoTimelineEntry is one STAGE log transition.
type ReportInfoTimelineEntry struct {
	At        string
	FromStage string
	ToStage   string
}

// ReportInfoAppeal is the single appeal row from report_info.json.
type ReportInfoAppeal struct {
	Status        string
	TargetRole    string
	Message       string
	ReviewMessage string
	ReviewedAt    string
}

// PDF represents the complete data structure for PDF generation
type PDF struct {
	Report    Report      `json:"report"`
	Judge     Judgment    `json:"judge"`
	Amendment []Amendment `json:"amendment,omitempty"`
}

// Comment is a flattened thread entry from comments.json (backend-api ZIPs).
type Comment struct {
	Message   string        `json:"message"`
	RepliedTo string        `json:"replied_to,omitempty"`
	By        CommentAuthor `json:"by"`
	CreatedAt string        `json:"created_at"`
}

// CommentAuthor identifies who posted a comment.
type CommentAuthor struct {
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display_name,omitempty"`
	Role        string `json:"role"`
}

// CommentsFile is the plaintext comments.json envelope.
type CommentsFile struct {
	SchemaVersion int       `json:"schema_version"`
	Comments      []Comment `json:"comments"`
}
