package pdfgen

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"ravro_dcrpt/pkg/models"
)

func solutionText(report *models.Report) string {
	if strings.TrimSpace(report.Solution) != "" {
		return report.Solution
	}
	return report.Description
}

func submittedAt(report *models.Report) string {
	if strings.TrimSpace(report.CreatedAt) != "" {
		return formatDateTimeFA(report.CreatedAt)
	}
	return formatDateTimeFA(report.SubmissionDate)
}

func isHTTPURL(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return false
	}
	u, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func safeExternalHref(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if !isHTTPURL(trimmed) {
		return ""
	}
	return trimmed
}

func renderNamedLink(title, href string) string {
	name := esc(title)
	link := safeExternalHref(href)
	if link == "" {
		return name
	}
	return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`, esc(link), name)
}

func isAvatarCode(value string) bool {
	v := strings.TrimSpace(value)
	if v == "" {
		return false
	}
	if isHTTPURL(v) {
		return false
	}
	if strings.ContainsAny(v, ":/\\") || strings.HasPrefix(v, ".") {
		return false
	}
	return true
}

func inferAvatarAPIBase(urls ...string) string {
	for _, raw := range urls {
		u, err := url.Parse(strings.TrimSpace(raw))
		if err != nil || u.Hostname() == "" {
			continue
		}
		host := strings.ToLower(u.Hostname())
		scheme := u.Scheme
		if scheme == "" {
			scheme = "https"
		}
		if host == "localhost" || host == "127.0.0.1" {
			return "http://localhost:4000/api/v1"
		}
		host = strings.TrimPrefix(host, "www.")
		if !strings.HasPrefix(host, "api.") {
			host = "api." + host
		}
		return scheme + "://" + host + "/api/v1"
	}
	return "https://api.ravro.ir/api/v1"
}

func publicAvatarSrc(code, apiBase string) string {
	trimmed := strings.TrimSpace(code)
	if isHTTPURL(trimmed) {
		return trimmed
	}
	if !isAvatarCode(trimmed) {
		return ""
	}
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if base == "" {
		return ""
	}
	return base + "/public/attachments/avatar?code=" + url.QueryEscape(trimmed)
}

func renderPartyAvatar(title, avatar, apiBase string) string {
	initial := esc(avatarInitial(title))
	src := publicAvatarSrc(avatar, apiBase)
	if src == "" {
		return fmt.Sprintf(`<span class="avatar">%s</span>`, initial)
	}
	return fmt.Sprintf(`<span class="avatar-wrap"><img class="avatar-img" src="%s" alt="" data-initial="%s" onerror="var s=document.createElement('span');s.className='avatar';s.textContent=this.getAttribute('data-initial')||'?';this.replaceWith(s)" /></span>`,
		esc(src), initial)
}

func renderPartyChip(title, href, avatar, apiBase string) string {
	name := strings.TrimSpace(title)
	if name == "" {
		return esc("-")
	}
	avatarHTML := renderPartyAvatar(name, avatar, apiBase)
	link := safeExternalHref(href)
	if link == "" {
		return fmt.Sprintf(`<span class="party-chip">%s<span class="party-name">%s</span></span>`, avatarHTML, esc(name))
	}
	return fmt.Sprintf(`<span class="party-chip"><a class="party-link" href="%s" target="_blank" rel="noopener noreferrer">%s<span class="party-name">%s</span></a></span>`,
		esc(link), avatarHTML, esc(name))
}

func resolveDisplayParties(report *models.Report) (hunterHTML, companyHTML, targetHTML string, hasCompany bool) {
	exported := report.ReportInfo.Export
	hunterTitle := strings.TrimSpace(report.HunterUsername)
	companyTitle := strings.TrimSpace(report.CompanyUsername)
	targetTitle := reportTarget(report)
	hunterURL, companyURL, targetURL := "", "", ""
	hunterAvatar, companyAvatar := "", ""

	if exported.Present {
		if strings.TrimSpace(exported.Bughunter.Title) != "" {
			hunterTitle = exported.Bughunter.Title
		}
		hunterURL = exported.Bughunter.URL
		hunterAvatar = exported.Bughunter.Avatar
		if strings.TrimSpace(exported.Company.Title) != "" {
			companyTitle = exported.Company.Title
		}
		companyURL = exported.Company.URL
		companyAvatar = exported.Company.Avatar
		if strings.TrimSpace(exported.Target.Title) != "" {
			targetTitle = exported.Target.Title
		}
		targetURL = exported.Target.URL
	}

	apiBase := inferAvatarAPIBase(hunterURL, companyURL, targetURL)
	hunterHTML = renderPartyChip(displayOrDash(hunterTitle), hunterURL, hunterAvatar, apiBase)
	if strings.TrimSpace(companyTitle) != "" {
		hasCompany = true
		companyHTML = renderPartyChip(companyTitle, companyURL, companyAvatar, apiBase)
	}
	targetHTML = renderNamedLink(displayOrDash(targetTitle), targetURL)
	return hunterHTML, companyHTML, targetHTML, hasCompany
}

func renderInfoCell(label, valueHTML, extraClass string) string {
	class := "info-cell"
	if extraClass != "" {
		class += " " + extraClass
	}
	return fmt.Sprintf(`<div class="%s">
        <div class="cell-label">%s</div>
        <div class="cell-value">%s</div>
      </div>`, class, esc(label), valueHTML)
}

func formatBountyAmount(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	n, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return ""
	}
	intPart := int64(n)
	formatted := formatIntThousands(intPart)
	if n != float64(intPart) {
		formatted = formatFloatThousands(n)
	}
	return formatted + " " + labelBountyCurrency
}

func formatIntThousands(n int64) string {
	s := strconv.FormatInt(n, 10)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	if s != "" {
		parts = append([]string{s}, parts...)
	}
	out := strings.Join(parts, ",")
	if neg {
		return "-" + out
	}
	return out
}

func formatFloatThousands(n float64) string {
	s := strconv.FormatFloat(n, 'f', -1, 64)
	parts := strings.SplitN(s, ".", 2)
	intPart, _ := strconv.ParseInt(parts[0], 10, 64)
	out := formatIntThousands(intPart)
	if len(parts) == 2 {
		out += "." + parts[1]
	}
	return out
}

func renderExportInfoCells(exported models.ReportInfoExport) string {
	if !exported.Present {
		return ""
	}
	var cells []string
	if exported.UpdatedAt != "" {
		cells = append(cells, renderInfoCell(labelUpdatedAt, esc(formatDateTimeFA(exported.UpdatedAt)), ""))
	}
	if v := formatBountyAmount(exported.Bounty.Judge); v != "" {
		cells = append(cells, renderInfoCell(labelBountyJudge, esc(v), ""))
	}
	if v := formatBountyAmount(exported.Bounty.Company); v != "" {
		cells = append(cells, renderInfoCell(labelBountyCompany, esc(v), ""))
	}
	if v := formatBountyAmount(exported.Bounty.Final); v != "" {
		cells = append(cells, renderInfoCell(labelBountyFinal, esc(v), ""))
	}
	if exported.Bounty.PaidAt != "" {
		cells = append(cells, renderInfoCell(labelBountyPaidAt, esc(formatDateTimeFA(exported.Bounty.PaidAt)), ""))
	}
	if exported.DuplicateOfReportSlug != "" {
		cells = append(cells, renderInfoCell(labelDuplicateOf, `<span class="mono">`+esc(exported.DuplicateOfReportSlug)+`</span>`, ""))
	}
	if exported.Target.Version != "" {
		cells = append(cells, renderInfoCell(labelTargetVersion, esc(exported.Target.Version), ""))
	}
	return strings.Join(cells, "")
}

func renderRulesSection(rules string) string {
	if strings.TrimSpace(rules) == "" {
		return ""
	}
	return fmt.Sprintf(`<div class="section">
    <div class="section-head">
      <div class="section-icon">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
      </div>
      <span class="section-title">%s</span>
    </div>
    <div class="section-body">
      <div class="text-block md-content">%s</div>
    </div>
  </div>`, esc(labelRulesTitle), formatMarkdownHTML(rules))
}

func renderTimelineSection(timeline []models.ReportInfoTimelineEntry) string {
	if len(timeline) == 0 {
		return ""
	}
	var items strings.Builder
	for _, entry := range timeline {
		timeLabel := ""
		if strings.TrimSpace(entry.At) != "" {
			timeLabel = formatDateTimeFA(entry.At)
		}
		fromLabel := statusLabelFA(entry.FromStage)
		toLabel := statusLabelFA(entry.ToStage)
		if timeLabel == "" && fromLabel == "" && toLabel == "" {
			continue
		}
		arrow := ""
		if fromLabel != "" && toLabel != "" {
			arrow = " → "
		}
		timeHTML := ""
		if timeLabel != "" {
			timeHTML = fmt.Sprintf(`<span class="timeline-time">%s</span>`, esc(timeLabel))
		}
		fmt.Fprintf(&items, `<div class="timeline-row">
        %s
        <span class="timeline-stages">%s%s%s</span>
      </div>`, timeHTML, esc(fromLabel), arrow, esc(toLabel))
	}
	if items.Len() == 0 {
		return ""
	}
	return fmt.Sprintf(`<div class="section">
    <div class="section-head">
      <div class="section-icon" style="background:rgba(56,189,248,.1);color:#38bdf8;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
      </div>
      <span class="section-title">%s</span>
    </div>
    <div class="section-body">
      <div class="timeline-list">%s</div>
    </div>
  </div>`, esc(labelTimelineTitle), items.String())
}

func appealTypeLabel(role string) string {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case "COMPANY":
		return labelAppealTypeCompany
	case "JUDGE":
		return labelAppealTypeJudge
	default:
		return strings.TrimSpace(role)
	}
}

func appealStatusLabel(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "PENDING":
		return labelAppealPending
	case "ACCEPTED":
		return labelAppealAccepted
	case "REJECTED":
		return labelAppealRejected
	default:
		return strings.TrimSpace(status)
	}
}

func renderAppealSection(appeal *models.ReportInfoAppeal) string {
	if appeal == nil {
		return ""
	}
	var rows []string
	if typeLabel := appealTypeLabel(appeal.TargetRole); typeLabel != "" {
		rows = append(rows, renderInfoCell(labelAppealType, esc(typeLabel), ""))
	}
	if strings.TrimSpace(appeal.Message) != "" {
		rows = append(rows, renderInfoCell(labelAppealMessage, `<span class="plain-text">`+esc(appeal.Message)+`</span>`, "span-3"))
	}
	if strings.TrimSpace(appeal.ReviewMessage) != "" {
		rows = append(rows, renderInfoCell(labelAppealReview, `<span class="plain-text">`+esc(appeal.ReviewMessage)+`</span>`, "span-3"))
	}
	if strings.TrimSpace(appeal.ReviewedAt) != "" {
		rows = append(rows, renderInfoCell(labelAppealReviewedAt, esc(formatDateTimeFA(appeal.ReviewedAt)), ""))
	}
	status := appealStatusLabel(appeal.Status)
	count := ""
	if status != "" {
		count = fmt.Sprintf(`<span class="section-count">%s</span>`, esc(status))
	}
	return fmt.Sprintf(`<div class="section">
    <div class="section-head">
      <div class="section-icon" style="background:rgba(244,63,94,.1);color:#fb7185;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"/></svg>
      </div>
      <span class="section-title">%s</span>
      %s
    </div>
    <div class="section-body">
      <div class="appeal-grid">%s</div>
    </div>
  </div>`, esc(labelAppealTitle), count, strings.Join(rows, ""))
}
