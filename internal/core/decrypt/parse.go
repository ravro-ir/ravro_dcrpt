package decrypt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ravro_dcrpt/pkg/models"
)

// DetectSchema returns SchemaVersionV2 when an envelope is present, otherwise V1.
func DetectSchema(raw []byte) int {
	var peek map[string]json.RawMessage
	if err := json.Unmarshal(raw, &peek); err != nil {
		return models.SchemaVersionV1
	}
	if _, ok := peek["schemaVersion"]; ok {
		return models.SchemaVersionV2
	}
	if _, ok := peek["schema_version"]; ok {
		return models.SchemaVersionV2
	}
	if _, ok := peek["payload"]; ok {
		return models.SchemaVersionV2
	}
	return models.SchemaVersionV1
}

// ParseReportJSON decrypts-to-JSON report ciphertext into the canonical Report model.
func ParseReportJSON(raw []byte) (*models.Report, int, error) {
	version := DetectSchema(raw)
	switch version {
	case models.SchemaVersionV2:
		var env reportV2Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, version, fmt.Errorf("failed to parse v2 report JSON: %w", err)
		}
		report := normalizeReportV2(env)
		return report, models.SchemaVersionV2, nil
	default:
		report, err := parseReportV1(raw)
		if err != nil {
			return nil, models.SchemaVersionV1, err
		}
		return report, models.SchemaVersionV1, nil
	}
}

// ParseAmendmentJSON parses amendment ciphertext (v1 flat or v2 envelope).
func ParseAmendmentJSON(raw []byte) (*models.Amendment, error) {
	switch DetectSchema(raw) {
	case models.SchemaVersionV2:
		var env amendmentV2Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("failed to parse v2 amendment JSON: %w", err)
		}
		a := normalizeAmendmentV2(env)
		return &a, nil
	default:
		var rawV1 amendmentV1Raw
		if err := json.Unmarshal(raw, &rawV1); err != nil {
			return nil, fmt.Errorf("failed to parse v1 amendment JSON: %w", err)
		}
		a := models.Amendment{
			ReportID:         rawV1.ReportID.String(),
			AttachmentsCount: rawV1.AttachmentsCount,
			Description:      rawV1.Description,
			HunterUsername:   rawV1.HunterUsername,
			CompanyUsername:  rawV1.CompanyUsername,
			SubmissionDate:   rawV1.SubmissionDate,
			Content:          rawV1.Content,
			CreatedAt:        rawV1.CreatedAt,
		}
		return &a, nil
	}
}

// ParseReportInfoJSON accepts management (flat hunter CVSS) and backend-api (nested cvss31/cvss40) shapes.
func ParseReportInfoJSON(raw []byte) (*models.InfoReport, error) {
	var wire reportInfoWire
	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, fmt.Errorf("failed to parse report_info.json: %w", err)
	}

	info := &models.InfoReport{
		SchemaVersion: wire.SchemaVersion,
		Tags:          make([]models.TagInfo, 0, len(wire.Tags)),
		Export:        parseExportBlock(wire.Export),
	}
	for _, tag := range wire.Tags {
		info.Tags = append(info.Tags, models.TagInfo{
			InfoTitle:       firstNonEmpty(tag.InfoTitle, tag.InfoTitleSnake),
			InfoDescription: firstNonEmpty(tag.InfoDescription, tag.InfoDescriptionSnake),
			InfoSolution:    firstNonEmpty(tag.InfoSolution, tag.InfoSolutionSnake),
			InfoMore:        firstNonEmpty(tag.InfoMore, tag.InfoMoreSnake),
		})
	}

	info.Details.Target = wire.Details.Target
	info.Details.CurrentStatus = wire.Details.CurrentStatus
	info.Details.Judges = wire.Details.Judges
	info.Details.Attachments = wire.Details.Attachments
	info.Details.Cvss.Hunter = pickRoleCVSS(wire.Details.Cvss, "hunter", "bughunter")
	info.Details.Cvss.Judge = pickRoleCVSS(wire.Details.Cvss, "judge")
	info.Details.Cvss.Company = pickRoleCVSS(wire.Details.Cvss, "company")
	info.Details.Cvss.Final = pickRoleCVSS(wire.Details.Cvss, "final")

	return info, nil
}

type commentsFileWire struct {
	SchemaVersion int               `json:"schema_version"`
	Comments      []json.RawMessage `json:"comments"`
}

type commentWire struct {
	Message      string          `json:"message"`
	Text         string          `json:"text"`
	Content      string          `json:"content"`
	RepliedTo    string          `json:"replied_to"`
	RepliedToAlt string          `json:"repliedTo"`
	Nickname     string          `json:"nickname"`
	Author       json.RawMessage `json:"author"`
	User         string          `json:"user"`
	Role         string          `json:"role"`
	AuthorRole   string          `json:"author_role"`
	By           json.RawMessage `json:"by"`
	CreatedAt    json.RawMessage `json:"created_at"`
	CreatedAtAlt json.RawMessage `json:"createdAt"`
	Time         string          `json:"time"`
}

type commentAuthorWire struct {
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display_name"`
	Title       string `json:"title"`
	Role        string `json:"role"`
}

// ParseCommentsJSON parses optional comments.json from backend-api ZIPs.
func ParseCommentsJSON(raw []byte) ([]models.Comment, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return nil, nil
	}
	schemaVersion := 1
	var items []json.RawMessage
	if strings.HasPrefix(trimmed, "[") {
		if err := json.Unmarshal([]byte(trimmed), &items); err != nil {
			return nil, fmt.Errorf("failed to parse comments.json: %w", err)
		}
	} else {
		var file commentsFileWire
		if err := json.Unmarshal([]byte(trimmed), &file); err != nil {
			return nil, fmt.Errorf("failed to parse comments.json: %w", err)
		}
		items = file.Comments
		if file.SchemaVersion > 0 {
			schemaVersion = file.SchemaVersion
		}
	}
	out := make([]models.Comment, 0, len(items))
	for _, item := range items {
		comment, ok := normalizeCommentJSON(item, schemaVersion)
		if ok {
			out = append(out, comment)
		}
	}
	return out, nil
}

func normalizeCommentJSON(raw json.RawMessage, schemaVersion int) (models.Comment, bool) {
	var wire commentWire
	if err := json.Unmarshal(raw, &wire); err != nil {
		return models.Comment{}, false
	}
	by, byName := parseCommentBy(wire.By)
	role := firstNonEmpty(wire.Role, wire.AuthorRole, by.Role, "BUGHUNTER")
	nickname := firstNonEmpty(
		wire.Nickname,
		jsonRawString(wire.Author),
		wire.User,
		by.Nickname,
		by.Title,
		by.DisplayName,
		byName,
	)
	if schemaVersion >= 2 && by.DisplayName != "" {
		nickname = by.DisplayName
	}
	message := firstNonEmpty(wire.Message, wire.Text, wire.Content)
	if strings.TrimSpace(nickname) == "" && strings.TrimSpace(message) == "" {
		return models.Comment{}, false
	}
	return models.Comment{
		Message:   message,
		RepliedTo: firstNonEmpty(wire.RepliedTo, wire.RepliedToAlt),
		By: models.CommentAuthor{
			Nickname:    nickname,
			DisplayName: by.DisplayName,
			Role:        role,
		},
		CreatedAt: firstNonEmpty(jsonRawString(wire.CreatedAt), jsonRawString(wire.CreatedAtAlt), wire.Time),
	}, true
}

func parseCommentBy(raw json.RawMessage) (commentAuthorWire, string) {
	if len(raw) == 0 || string(raw) == "null" {
		return commentAuthorWire{}, ""
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return commentAuthorWire{}, asString
	}
	var by commentAuthorWire
	_ = json.Unmarshal(raw, &by)
	return by, ""
}

func jsonRawString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return ""
}

// MergeReportInfo overlays ZIP metadata onto a decrypted report and fills CVSS/target gaps.
func MergeReportInfo(report *models.Report, info *models.InfoReport) {
	if report == nil || info == nil {
		return
	}
	seedTarget := report.ReportInfo.Details.Target
	seedHunter := report.ReportInfo.Details.Cvss.Hunter
	seedJudge := report.ReportInfo.Details.Cvss.Judge
	seedCompany := report.ReportInfo.Details.Cvss.Company
	seedFinal := report.ReportInfo.Details.Cvss.Final
	seedExport := report.ReportInfo.Export

	report.ReportInfo = *info

	if report.ReportInfo.Details.Target == "" {
		report.ReportInfo.Details.Target = firstNonEmpty(seedTarget, report.TargetName)
	}
	if report.ReportInfo.Details.Cvss.Hunter.Vector == "" {
		report.ReportInfo.Details.Cvss.Hunter.Vector = firstNonEmpty(seedHunter.Vector, report.Cvss)
	}
	if report.ReportInfo.Details.Cvss.Hunter.Score == "" {
		report.ReportInfo.Details.Cvss.Hunter.Score = seedHunter.Score
	}
	if report.ReportInfo.Details.Cvss.Judge.Vector == "" {
		report.ReportInfo.Details.Cvss.Judge = seedJudge
	}
	if report.ReportInfo.Details.Cvss.Company.Vector == "" && report.ReportInfo.Details.Cvss.Company.Score == "" {
		report.ReportInfo.Details.Cvss.Company = seedCompany
	}
	if report.ReportInfo.Details.Cvss.Final.Vector == "" && report.ReportInfo.Details.Cvss.Final.Score == "" {
		report.ReportInfo.Details.Cvss.Final = seedFinal
	}
	if !report.ReportInfo.Export.Present {
		report.ReportInfo.Export = seedExport
	}
}

func parseReportV1(raw []byte) (*models.Report, error) {
	var report models.Report
	if err := json.Unmarshal(raw, &report); err != nil {
		return nil, fmt.Errorf("failed to parse v1 report JSON: %w", err)
	}
	report.SchemaVersion = models.SchemaVersionV1
	if report.ReportInfo.Details.Target == "" {
		report.ReportInfo.Details.Target = report.TargetName
	}
	if report.ReportInfo.Details.Cvss.Hunter.Vector == "" {
		report.ReportInfo.Details.Cvss.Hunter.Vector = report.Cvss
	}
	if strings.TrimSpace(report.Solution) == "" {
		report.Solution = report.Description
	}
	return &report, nil
}

func normalizeReportV2(env reportV2Envelope) *models.Report {
	p := env.Payload
	report := &models.Report{
		Slug:            p.Slug,
		Title:           p.Title,
		Description:     p.Description,
		Solution:        firstNonEmpty(p.Solution, p.Description),
		Scenario:        p.Scenario,
		HunterUsername:  entityLabel(p.Bughunter),
		CompanyUsername: entityLabel(p.Company),
		Urls:            p.URLs,
		Ips:             p.IPs,
		SubmissionDate:  p.SubmissionDate,
		DateFrom:        p.DateFrom,
		DateTo:          p.DateTo,
		TargetName:      entityLabel(p.Target),
		Cvss:            p.CVSS,
		CreatedAt:       env.CreatedAt,
		SchemaVersion:   models.SchemaVersionV2,
	}
	if env.SchemaVersion != 0 {
		report.SchemaVersion = env.SchemaVersion
	}
	report.ReportInfo.Details.Target = entityTitle(p.Target)
	report.ReportInfo.Details.Cvss.Hunter.Vector = p.CVSS
	return report
}

func normalizeAmendmentV2(env amendmentV2Envelope) models.Amendment {
	a := models.Amendment{
		ReportID:    env.Payload.ReportSlug,
		Description: env.Payload.Explanation,
		FromStage:   env.Payload.FromStage,
		ToStage:     env.Payload.ToStage,
	}
	a.AttachmentsCount = len(env.Payload.Attachments)
	if !env.CreatedAt.IsZero() {
		a.CreatedAt = env.CreatedAt
	}
	return a
}

func entityLabel(e entityRef) string {
	return firstNonEmpty(e.Slug, e.Title)
}

func entityTitle(e entityRef) string {
	return firstNonEmpty(e.Title, e.Slug)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func pickRoleCVSS(raw json.RawMessage, keys ...string) models.CVSSScore {
	if len(raw) == 0 || string(raw) == "null" {
		return models.CVSSScore{}
	}
	var roles map[string]json.RawMessage
	if err := json.Unmarshal(raw, &roles); err != nil {
		return models.CVSSScore{}
	}
	for _, key := range keys {
		blob, ok := roles[key]
		if !ok || len(blob) == 0 || string(blob) == "null" {
			continue
		}
		if score := parseCVSSBlob(blob); score.Vector != "" || score.Score != "" || score.Rating != "" {
			return score
		}
	}
	return models.CVSSScore{}
}

func parseCVSSBlob(raw json.RawMessage) models.CVSSScore {
	var nested struct {
		CVSS40 *nullableScore `json:"cvss40"`
		CVSS31 *nullableScore `json:"cvss31"`
		Vector string         `json:"vector"`
		Score  string         `json:"score"`
		Rating string         `json:"rating"`
	}
	if err := json.Unmarshal(raw, &nested); err != nil {
		return models.CVSSScore{}
	}
	if nested.CVSS40 != nil {
		if s := nested.CVSS40.toScore(); s.Vector != "" || s.Score != "" {
			return s
		}
	}
	if nested.CVSS31 != nil {
		if s := nested.CVSS31.toScore(); s.Vector != "" || s.Score != "" {
			return s
		}
	}
	return models.CVSSScore{
		Vector: nested.Vector,
		Score:  nested.Score,
		Rating: nested.Rating,
	}
}

type nullableScore struct {
	Vector *string `json:"vector"`
	Score  *string `json:"score"`
}

func (n *nullableScore) toScore() models.CVSSScore {
	if n == nil {
		return models.CVSSScore{}
	}
	s := models.CVSSScore{}
	if n.Vector != nil {
		s.Vector = *n.Vector
	}
	if n.Score != nil {
		s.Score = *n.Score
	}
	return s
}

type reportV2Envelope struct {
	SchemaVersion int             `json:"schemaVersion"`
	CreatedAt     string          `json:"createdAt"`
	Payload       reportV2Payload `json:"payload"`
}

type reportV2Payload struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Solution         string    `json:"solution"`
	Scenario         string    `json:"scenario"`
	SubmissionDate   string    `json:"submission_date"`
	Slug             string    `json:"slug"`
	AttachmentsCount int       `json:"attachments_count"`
	IPs              string    `json:"ips"`
	URLs             string    `json:"urls"`
	CVSS             string    `json:"cvss"`
	DateFrom         string    `json:"date_from"`
	DateTo           string    `json:"date_to"`
	Bughunter        entityRef `json:"bughunter"`
	Company          entityRef `json:"company"`
	Workspace        entityRef `json:"workspace"`
	Target           entityRef `json:"target"`
}

type entityRef struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type amendmentV2Envelope struct {
	SchemaVersion int                `json:"schema_version"`
	CreatedAt     time.Time          `json:"created_at"`
	Payload       amendmentV2Payload `json:"payload"`
}

type amendmentV2Payload struct {
	Explanation string `json:"explanation"`
	Attachments []struct {
		AttachmentID string `json:"attachment_id"`
		ExtHash      string `json:"ext_hash"`
	} `json:"attachments"`
	ReportSlug string `json:"report_slug"`
	FromStage  string `json:"from_stage"`
	ToStage    string `json:"to_stage"`
}

type amendmentV1Raw struct {
	ReportID         flexibleID `json:"reportId"`
	AttachmentsCount int        `json:"attachmentsCount"`
	Description      string     `json:"description"`
	HunterUsername   string     `json:"hunterUsername"`
	CompanyUsername  string     `json:"companyUsername"`
	SubmissionDate   string     `json:"submissionDate"`
	Content          string     `json:"content"`
	CreatedAt        time.Time  `json:"created_at"`
}

type flexibleID string

func (id *flexibleID) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || len(b) == 0 {
		*id = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*id = flexibleID(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		*id = flexibleID(n.String())
		return nil
	}
	var i int
	if err := json.Unmarshal(b, &i); err == nil {
		*id = flexibleID(strconv.Itoa(i))
		return nil
	}
	return fmt.Errorf("invalid reportId: %s", string(b))
}

func (id flexibleID) String() string {
	return string(id)
}

type reportInfoWire struct {
	SchemaVersion int               `json:"schema_version"`
	Tags          []tagInfoWire     `json:"tags"`
	Details       reportDetailsWire `json:"details"`
	Export        json.RawMessage   `json:"export"`
}

type tagInfoWire struct {
	InfoTitle            string `json:"infoTitle"`
	InfoTitleSnake       string `json:"info_title"`
	InfoDescription      string `json:"infoDescription"`
	InfoDescriptionSnake string `json:"info_description"`
	InfoSolution         string `json:"infoSolution"`
	InfoSolutionSnake    string `json:"info_solution"`
	InfoMore             string `json:"infoMore"`
	InfoMoreSnake        string `json:"info_more"`
}

type reportDetailsWire struct {
	Target        string              `json:"target"`
	CurrentStatus string              `json:"currentStatus"`
	Cvss          json.RawMessage     `json:"cvss"`
	Attachments   []models.Attachment `json:"attachments"`
	Judges        []models.Judge      `json:"judges"`
}

func toSlashPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// IsReportPayloadFile is true for report/data.ravro (main ciphertext).
func IsReportPayloadFile(path string) bool {
	p := toSlashPath(path)
	return strings.Contains(p, "/report/") && strings.HasSuffix(p, "/data.ravro")
}

// IsJudgmentPayloadFile is true for judgment/data.ravro.
func IsJudgmentPayloadFile(path string) bool {
	p := toSlashPath(path)
	return strings.Contains(p, "/judgment/") && strings.HasSuffix(p, "/data.ravro")
}

// IsAmendmentPayloadFile is true for amendment-N or ammendment-N data.ravro / data.com.ravro.
func IsAmendmentPayloadFile(path string) bool {
	p := toSlashPath(path)
	if strings.Contains(p, "/report/") || strings.Contains(p, "/judgment/") {
		return false
	}
	inFolder := strings.Contains(p, "/amendment") || strings.Contains(p, "/ammendment")
	if !inFolder {
		return false
	}
	return strings.HasSuffix(p, "/data.ravro") || strings.HasSuffix(p, "/data.com.ravro")
}

// IsStructuredDataFile is true for ciphertext that is JSON payload, not an attachment.
func IsStructuredDataFile(path string) bool {
	p := toSlashPath(path)
	return strings.HasSuffix(p, "/data.ravro") || strings.HasSuffix(p, "/data.com.ravro")
}
