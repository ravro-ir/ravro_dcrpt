package decrypt

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"ravro_dcrpt/pkg/models"
)

func parseExportBlock(raw json.RawMessage) models.ReportInfoExport {
	if len(raw) == 0 || string(raw) == "null" {
		return models.ReportInfoExport{}
	}
	var wire struct {
		UpdatedAt             json.RawMessage      `json:"updated_at"`
		Bughunter             exportPartyWire      `json:"bughunter"`
		Company               exportPartyWire      `json:"company"`
		Target                exportTargetWire     `json:"target"`
		Bounty                exportBountyWire     `json:"bounty"`
		DuplicateOfReportSlug *string              `json:"duplicate_of_report_slug"`
		Timeline              []exportTimelineWire `json:"timeline"`
		Appeal                *exportAppealWire    `json:"appeal"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		return models.ReportInfoExport{}
	}
	out := models.ReportInfoExport{
		Present:               true,
		UpdatedAt:             jsonRawAsTimeString(wire.UpdatedAt),
		Bughunter:             wire.Bughunter.toParty(),
		Company:               wire.Company.toParty(),
		Target:                wire.Target.toTarget(),
		Bounty:                wire.Bounty.toBounty(),
		DuplicateOfReportSlug: derefString(wire.DuplicateOfReportSlug),
	}
	for _, item := range wire.Timeline {
		out.Timeline = append(out.Timeline, models.ReportInfoTimelineEntry{
			At:        jsonRawAsTimeString(item.At),
			FromStage: item.FromStage,
			ToStage:   item.ToStage,
		})
	}
	if wire.Appeal != nil {
		out.Appeal = &models.ReportInfoAppeal{
			Status:        wire.Appeal.Status,
			TargetRole:    wire.Appeal.TargetRole,
			Message:       wire.Appeal.Message,
			ReviewMessage: derefString(wire.Appeal.ReviewMessage),
			ReviewedAt:    jsonRawAsTimeString(wire.Appeal.ReviewedAt),
		}
	}
	return out
}

type exportPartyWire struct {
	Title  string  `json:"title"`
	Slug   string  `json:"slug"`
	URL    *string `json:"url"`
	Avatar *string `json:"avatar"`
}

func (p exportPartyWire) toParty() models.ReportInfoParty {
	return models.ReportInfoParty{
		Title:  p.Title,
		Slug:   p.Slug,
		URL:    derefString(p.URL),
		Avatar: derefString(p.Avatar),
	}
}

type exportTargetWire struct {
	Title   string          `json:"title"`
	Slug    string          `json:"slug"`
	URL     *string         `json:"url"`
	Version json.RawMessage `json:"version"`
	Rules   *string         `json:"rules"`
}

func (t exportTargetWire) toTarget() models.ReportInfoExportTarget {
	return models.ReportInfoExportTarget{
		Title:   t.Title,
		Slug:    t.Slug,
		URL:     derefString(t.URL),
		Version: jsonRawAsPlainString(t.Version),
		Rules:   derefString(t.Rules),
	}
}

type exportBountyWire struct {
	Judge   *string         `json:"judge"`
	Company *string         `json:"company"`
	Final   *string         `json:"final"`
	PaidAt  json.RawMessage `json:"paid_at"`
}

func (b exportBountyWire) toBounty() models.ReportInfoExportBounty {
	return models.ReportInfoExportBounty{
		Judge:   derefString(b.Judge),
		Company: derefString(b.Company),
		Final:   derefString(b.Final),
		PaidAt:  jsonRawAsTimeString(b.PaidAt),
	}
}

type exportTimelineWire struct {
	At        json.RawMessage `json:"at"`
	FromStage string          `json:"from_stage"`
	ToStage   string          `json:"to_stage"`
}

type exportAppealWire struct {
	Status        string          `json:"status"`
	TargetRole    string          `json:"target_role"`
	Message       string          `json:"message"`
	ReviewMessage *string         `json:"review_message"`
	ReviewedAt    json.RawMessage `json:"reviewed_at"`
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func jsonRawAsTimeString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strings.TrimSpace(s)
	}
	var t time.Time
	if err := json.Unmarshal(raw, &t); err == nil && !t.IsZero() {
		return t.UTC().Format(time.RFC3339)
	}
	return ""
}

func jsonRawAsPlainString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strings.TrimSpace(s)
	}
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		return n.String()
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return strings.Trim(string(raw), `"`)
}
