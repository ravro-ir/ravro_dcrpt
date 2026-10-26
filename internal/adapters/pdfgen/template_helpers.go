package pdfgen

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	ptime "github.com/yaa110/go-persian-calendar"

	"ravro_dcrpt/pkg/models"
)

type severityInfo struct {
	level    string
	score    float64
	cssClass string
}

type statusPillInfo struct {
	label    string
	cssClass string
	pulse    bool
}

type cvssResolved struct {
	versionLabel string
	score        *float64
	severity     *severityInfo
	vector       string
}

const sectionChevronSVG = `<svg class="section-chevron" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>`

func collapsibleSection(headHTML, bodyHTML string) string {
	return fmt.Sprintf(`<details class="section" open>
    <summary class="section-head">%s%s</summary>
    <div class="section-body">%s</div>
  </details>`, headHTML, sectionChevronSVG, bodyHTML)
}

func buildReportHTML(report *models.Report, judgment *models.Judgment, amendments []models.Amendment, outputDir string) string {
	slug := report.Slug
	title := report.Title
	if strings.TrimSpace(title) == "" {
		title = slug
	}

	severity := heroSeverity(report, judgment)
	severityLabel := severityLabelFA(severity.level)
	createdAt := submittedAt(report)
	generationDate := formatDateTimeFA(time.Now().Format(time.RFC3339))
	hunterHTML, companyHTML, targetHTML, hasCompany := resolveDisplayParties(report)
	activityRange := activityRangeFA(report)
	status := statusPill(report.ReportInfo.Details.CurrentStatus)
	urls := renderURLs(report.Urls)
	companyMeta := ""
	companyCell := ""
	if hasCompany {
		companyMeta = fmt.Sprintf(`<span class="meta-item">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="7" width="18" height="13" rx="2"/><path d="M8 7V5a4 4 0 018 0v2"/></svg>
              %s
            </span>`, companyHTML)
		companyCell = fmt.Sprintf(`<div class="info-cell">
        <div class="cell-label">%s</div>
        <div class="cell-value">%s</div>
      </div>`, esc(labelCompany), companyHTML)
	}

	heroSeverityBlock := ""
	if severity.level != "NONE" {
		heroSeverityBlock = fmt.Sprintf(`<div class="severity-block %s">
          <span class="severity-score">%.1f</span>
          <span class="severity-badge %s">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            %s
          </span>
        </div>`, severity.cssClass, severity.score, severity.cssClass, esc(severityLabel))
	}

	statusDot := `<span class="status-dot" style="animation:none;"></span>`
	pulseClass := ""
	if status.pulse {
		statusDot = `<span class="status-dot"></span>`
		pulseClass = " pulse"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="fa" dir="rtl" data-theme="light">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>%s</title>
  <style>%s</style>
</head>
<body>
  <div class="topbar no-print">
    <div class="topbar-left">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
      %s — %s
    </div>
    <div class="topbar-actions">
      <button type="button" class="btn" id="theme-toggle" aria-label="%s">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/></svg>
      </button>
      <button type="button" class="btn" onclick="window.print()">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 01-2-2v-5a2 2 0 012-2h16a2 2 0 012 2v5a2 2 0 01-2 2h-2"/><rect x="6" y="14" width="12" height="8"/></svg>
        %s
      </button>
    </div>
  </div>

  <div class="card">
    <div class="hero">
      <div class="hero-top">
        <div>
          <div class="hero-id">%s</div>
          <h1 class="hero-title">%s</h1>
          <div class="meta-row" style="margin-top:12px;">
            <span class="meta-item">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              %s: %s
            </span>
            <span class="meta-item">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              %s
            </span>
            %s
            <span class="meta-item">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              %s
            </span>
          </div>
        </div>
        <div class="hero-right">
          %s
          <span class="status-pill %s%s">
            %s
            %s
          </span>
        </div>
      </div>
    </div>

    <div class="info-grid">
      <div class="info-cell">
        <div class="cell-label">%s</div>
        <div class="cell-value mono" style="direction:ltr;text-align: right;">%s</div>
      </div>
      <div class="info-cell">
        <div class="cell-label">%s</div>
        <div class="cell-value">%s</div>
      </div>
      %s
      <div class="info-cell">
        <div class="cell-label">%s</div>
        <div class="cell-value">%s</div>
      </div>
      <div class="info-cell">
        <div class="cell-label">%s</div>
        <div class="cell-value">%s</div>
      </div>
      <div class="info-cell">
        <div class="cell-label">%s</div>
        <div class="cell-value">%s</div>
      </div>
      <div class="info-cell">
        <div class="cell-label">%s</div>
        <div class="cell-value mono" style="direction:ltr;text-align: right;">%s</div>
      </div>
      %s
      <div class="info-cell span-3">
        <div class="cell-label">%s</div>
        <div class="cell-value">%s</div>
      </div>
    </div>

    %s
    %s

    %s

    %s

    %s

    %s
    %s
    %s
  </div>

  <div class="report-footer no-print">
    <div class="footer-brand">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
      %s
    </div>
    <div style="font-size:11px;color:var(--text-tertiary);">%s</div>
  </div>
  %s
</body>
</html>`,
		esc(title),
		reportHTMLStyles,
		esc(labelTopbarTitle), esc(slug),
		esc(labelThemeDark),
		esc(labelPrint),
		esc(slug),
		esc(title),
		esc(labelSubmittedAt), esc(createdAt),
		hunterHTML,
		companyMeta,
		targetHTML,
		heroSeverityBlock,
		esc(status.cssClass), pulseClass,
		statusDot,
		esc(status.label),
		esc(labelReportID), esc(slug),
		esc(labelHunter), hunterHTML,
		companyCell,
		esc(labelSubmittedAt), esc(createdAt),
		esc(labelActivityRange), esc(activityRange),
		esc(labelTarget), targetHTML,
		esc(labelIPs), esc(formatList(report.Ips)),
		renderExportInfoCells(report.ReportInfo.Export),
		esc(labelURLs), urls,
		renderTimelineSection(report.ReportInfo.Export.Timeline),
		renderAppealSection(report.ReportInfo.Export.Appeal),
		collapsibleSection(`<div class="section-icon">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
        </div>
        <span class="section-title">`+esc(labelCVSSSection)+`</span>`, renderCVSSTable(report, judgment)),
		collapsibleSection(`<div class="section-icon">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        </div>
        <span class="section-title">`+esc(labelScenario)+`</span>`, `<div class="text-block md-content">`+formatMarkdownHTML(report.Scenario)+`</div>`),
		collapsibleSection(`<div class="section-icon" style="background:rgba(34,197,94,.1);color:#4ade80;">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        </div>
        <span class="section-title">`+esc(labelSolution)+`</span>`, `<div class="text-block md-content">`+formatMarkdownHTML(solutionText(report))+`</div>`),
		renderAmendments(amendments, outputDir),
		renderAttachments(resolveReportAttachments(outputDir, report.ReportInfo.Details.Attachments)),
		renderComments(report.Comments),
		esc(labelFooterBrand),
		esc(generationDate),
		themeToggleScript(),
	)
}

func reportTarget(report *models.Report) string {
	if report.ReportInfo.Details.Target != "" {
		return report.ReportInfo.Details.Target
	}
	return report.TargetName
}

func displayOrDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func formatList(items string) string {
	parts := splitList(items)
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, ", ")
}

func splitList(items string) []string {
	if strings.TrimSpace(items) == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(items, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func activityRangeFA(report *models.Report) string {
	from := formatDateFA(report.DateFrom)
	to := formatDateFA(report.DateTo)
	if from == "-" || to == "-" {
		return "-"
	}
	return from + " — " + to
}

func formatDateFA(dateStr string) string {
	shamsi := convertToShamsi(dateStr)
	if strings.TrimSpace(shamsi) == "" {
		return "-"
	}
	return shamsi
}

func formatDateTimeFA(dateStr string) string {
	if strings.TrimSpace(dateStr) == "" {
		return "-"
	}
	t, ok := parseFlexibleTime(dateStr)
	if !ok {
		if shamsi := convertToShamsi(dateStr); shamsi != "" {
			return shamsi
		}
		return "-"
	}
	pt := ptime.New(t.In(ptime.Iran()))
	return pt.Format("yyyy/MM/dd HH:mm")
}

func parseFlexibleTime(dateStr string) (time.Time, bool) {
	dateStr = strings.TrimSpace(dateStr)
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func convertToShamsi(dateStr string) string {
	if strings.TrimSpace(dateStr) == "" {
		return ""
	}
	t, ok := parseFlexibleTime(dateStr)
	if !ok {
		parts := strings.Split(strings.TrimSpace(dateStr), " ")
		dateFields := strings.Split(parts[0], "-")
		if len(dateFields) != 3 {
			return dateStr
		}
		year, err1 := strconv.Atoi(dateFields[0])
		month, err2 := strconv.Atoi(dateFields[1])
		day, err3 := strconv.Atoi(dateFields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return dateStr
		}
		t = time.Date(year, time.Month(month), day, 12, 0, 0, 0, ptime.Iran())
	}
	pt := ptime.New(t.In(ptime.Iran()))
	return pt.Format("yyyy/MM/dd")
}

func formatMarkdownHTML(mdText string) string {
	if strings.TrimSpace(mdText) == "" {
		return `<span class="null-value">-</span>`
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(mdText))
	opts := mdhtml.RendererOptions{Flags: mdhtml.CommonFlags | mdhtml.HrefTargetBlank}
	out := strings.TrimSpace(string(markdown.Render(doc, mdhtml.NewRenderer(opts))))
	if out == "" {
		return `<span class="null-value">-</span>`
	}
	return out
}

func renderURLs(raw string) string {
	urls := splitList(raw)
	if len(urls) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(urls))
	for _, u := range urls {
		parts = append(parts, fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`, esc(u), esc(u)))
	}
	return strings.Join(parts, "<br />")
}

func parseScore(value string) (float64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}

func severityCSSClass(level string) string {
	switch level {
	case "CRITICAL":
		return "sev-critical"
	case "HIGH":
		return "sev-high"
	case "MEDIUM":
		return "sev-medium"
	case "LOW":
		return "sev-low"
	default:
		return "sev-info"
	}
}

func severityFromScore(score float64) severityInfo {
	switch {
	case score >= 9.0:
		return severityInfo{level: "CRITICAL", score: score, cssClass: "sev-critical"}
	case score >= 7.0:
		return severityInfo{level: "HIGH", score: score, cssClass: "sev-high"}
	case score >= 4.0:
		return severityInfo{level: "MEDIUM", score: score, cssClass: "sev-medium"}
	case score > 0:
		return severityInfo{level: "LOW", score: score, cssClass: "sev-low"}
	default:
		return severityInfo{level: "NONE", score: 0, cssClass: "sev-info"}
	}
}

func normalizeSeverityLevel(level string) string {
	n := strings.ToUpper(strings.TrimSpace(level))
	switch n {
	case "CRITICAL", "HIGH", "MEDIUM", "LOW", "NONE":
		return n
	case "بحرانی":
		return "CRITICAL"
	case "بالا":
		return "HIGH"
	case "متوسط":
		return "MEDIUM"
	case "پایین":
		return "LOW"
	default:
		return ""
	}
}

func severityFromScoreAndLevel(score float64, level string) severityInfo {
	if n := normalizeSeverityLevel(level); n != "" {
		return severityInfo{level: n, score: score, cssClass: severityCSSClass(n)}
	}
	return severityFromScore(score)
}

func detectCVSSVersion(vector string) string {
	if strings.TrimSpace(vector) == "" {
		return "—"
	}
	if strings.Contains(vector, "CVSS:4.0") {
		return "CVSS 4.0"
	}
	return "CVSS 3.1"
}

func resolveCVSS(score models.CVSSScore, fallbackVector string) *cvssResolved {
	vector := strings.TrimSpace(score.Vector)
	if vector == "" {
		vector = strings.TrimSpace(fallbackVector)
	}
	parsed, hasScore := parseScore(score.Score)
	if !hasScore && vector == "" {
		return nil
	}
	out := &cvssResolved{
		versionLabel: detectCVSSVersion(vector),
		vector:       vector,
	}
	if hasScore {
		sev := severityFromScoreAndLevel(parsed, score.Rating)
		out.score = &parsed
		out.severity = &sev
		return out
	}
	return out
}

func heroSeverity(report *models.Report, judgment *models.Judgment) severityInfo {
	candidates := []string{}
	if judgment != nil {
		candidates = append(candidates, judgment.Cvss.Value)
	}
	candidates = append(candidates,
		report.ReportInfo.Details.Cvss.Judge.Score,
		report.ReportInfo.Details.Cvss.Hunter.Score,
	)
	for _, c := range candidates {
		if score, ok := parseScore(c); ok {
			return severityFromScore(score)
		}
	}
	return severityInfo{level: "NONE", score: 0, cssClass: "sev-info"}
}

func statusPill(currentStatus string) statusPillInfo {
	raw := strings.TrimSpace(currentStatus)
	if raw == "" {
		raw = "PENDING"
	}
	status := strings.ToUpper(raw)
	label := statusLabelFA(status)
	info := statusPillInfo{label: label, cssClass: "status-pending", pulse: false}

	switch status {
	case "COMPANY", "RAVRO", "RAVRO_REVISION_REQUEST", "RAVRO_REVISION_REVIEW", "RAVRO_OPINION":
		info.cssClass = "status-company"
		info.pulse = status == "COMPANY"
	case "JUDGE", "JUDGEMENT", "JUDGEMENT_REJECT", "JUDGEMENT_REVIEW":
		info.cssClass = "status-judge"
	case "RESOLVED", "ACCEPTED", "DISCLOSED", "DUPLICATE", "INFORMATIVE":
		info.cssClass = "status-resolved"
	case "REJECTED", "REJECTED_SPAM", "REJECTED_INCOMPLETE", "REJECTED_OUT_OF_SCOPE",
		"INITIAL_REJECT", "INITIAL_REJECT_SPAM", "JUDGEMENT_REJECT_CLOSE", "JUDGEMENT_REJECT_INCOMPLETE":
		info.cssClass = "status-rejected"
	}
	return info
}

func renderCVSSVersionRow(roleLabel string, resolved *cvssResolved, opts struct {
	isFinal    bool
	nullVector string
	nullScore  string
	showRole   bool
}) string {
	rowStyle := ""
	roleStyle := ""
	if opts.isFinal {
		rowStyle = ` style="background:rgba(99,102,241,.04);"`
		roleStyle = ` style="color:var(--text-primary);font-weight:700;"`
	}
	roleCell := ""
	if opts.showRole {
		roleCell = esc(roleLabel)
	}
	versionCell := "—"
	scoreCell := fmt.Sprintf(`<td class="null-value">%s</td>`, esc(opts.nullScore))
	severityCell := `<td class="null-value">—</td>`
	vectorCell := fmt.Sprintf(`<td class="null-value">%s</td>`, esc(opts.nullVector))
	if resolved != nil {
		if resolved.versionLabel != "" {
			versionCell = resolved.versionLabel
		}
		if resolved.score != nil && resolved.severity != nil {
			scoreCell = fmt.Sprintf(`<td class="cvss-score-cell %s" style="color:var(--sev-%s-text);">%.1f</td>`,
				resolved.severity.cssClass, strings.ToLower(resolved.severity.level), *resolved.score)
			severityCell = fmt.Sprintf(`<td><span class="severity-badge %s">%s</span></td>`,
				resolved.severity.cssClass, esc(severityLabelFA(resolved.severity.level)))
		}
		if resolved.vector != "" {
			vectorCell = fmt.Sprintf(`<td class="cvss-vector">%s</td>`, esc(resolved.vector))
		}
	}
	return fmt.Sprintf(`<tr%s>
    <td class="role-label"%s>%s</td>
    <td>%s</td>
    %s
    %s
    %s
  </tr>`, rowStyle, roleStyle, roleCell, esc(versionCell), scoreCell, severityCell, vectorCell)
}

func renderCVSSTable(report *models.Report, judgment *models.Judgment) string {
	hunter := resolveCVSS(report.ReportInfo.Details.Cvss.Hunter, report.Cvss)
	judgeScore := report.ReportInfo.Details.Cvss.Judge
	if judgment != nil {
		if judgeScore.Vector == "" {
			judgeScore.Vector = judgment.Cvss.Value
		}
		if _, ok := parseScore(judgeScore.Score); !ok {
			if _, ok := parseScore(judgment.Cvss.Rating); ok {
				judgeScore.Score = judgment.Cvss.Rating
			} else if judgeScore.Rating == "" {
				judgeScore.Rating = judgment.Cvss.Rating
			}
		}
	}
	judge := resolveCVSS(judgeScore, "")
	company := resolveCVSS(report.ReportInfo.Details.Cvss.Company, "")
	final := resolveCVSS(report.ReportInfo.Details.Cvss.Final, "")

	rows := renderCVSSVersionRow(labelRoleHunter, hunter, struct {
		isFinal    bool
		nullVector string
		nullScore  string
		showRole   bool
	}{nullVector: labelCVSSNullVector, nullScore: labelCVSSNullScore, showRole: true})
	rows += renderCVSSVersionRow(labelRoleJudge, judge, struct {
		isFinal    bool
		nullVector string
		nullScore  string
		showRole   bool
	}{nullVector: labelCVSSNullVector, nullScore: labelCVSSNullScore, showRole: true})
	rows += renderCVSSVersionRow(labelRoleCompany, company, struct {
		isFinal    bool
		nullVector string
		nullScore  string
		showRole   bool
	}{nullVector: labelCVSSNullVector, nullScore: labelCVSSNullScore, showRole: true})
	rows += renderCVSSVersionRow(labelRoleFinal, final, struct {
		isFinal    bool
		nullVector string
		nullScore  string
		showRole   bool
	}{isFinal: true, nullVector: labelCVSSFinalNull, nullScore: labelCVSSNullScore, showRole: true})

	return fmt.Sprintf(`<table class="cvss-table">
    <thead>
      <tr>
        <th>%s</th>
        <th>%s</th>
        <th>%s</th>
        <th>%s</th>
        <th>%s</th>
      </tr>
    </thead>
    <tbody>%s</tbody>
  </table>`, esc(labelCVSSRoleCol), esc(labelCVSSVersionCol), esc(labelCVSSScoreCol), esc(labelCVSSSeverityCol), esc(labelCVSSVectorCol), rows)
}

func fileIconMeta(filename string) (label, cssClass string) {
	ext := strings.ToLower(strings.TrimPrefix(path.Ext(filename), "."))
	switch ext {
	case "pdf":
		return "PDF", "file-icon-pdf"
	case "jpg", "jpeg":
		return "JPG", "file-icon-img"
	case "png":
		return "PNG", "file-icon-img"
	case "zip", "rar":
		return "ZIP", "file-icon-zip"
	case "mp4":
		return "MP4", "file-icon-vid"
	default:
		return "FILE", "file-icon-gen"
	}
}

func attachmentHref(relPath string) string {
	return encodeRelPath(relPath)
}

func renderAttachmentCards(attachments []resolvedAttachment) string {
	if len(attachments) == 0 {
		return `<div class="null-value">-</div>`
	}
	var b strings.Builder
	for _, att := range attachments {
		label, css := fileIconMeta(att.display)
		href := attachmentHref(att.relPath)
		fmt.Fprintf(&b, `<a href="%s" class="attachment-card">
              <div class="file-icon %s">%s</div>
              <div class="att-meta">
                <div class="att-name">%s</div>
              </div>
            </a>`, esc(href), css, esc(label), esc(att.display))
	}
	return b.String()
}

func renderAttachments(attachments []resolvedAttachment) string {
	countLabel := strings.ReplaceAll(labelAttachmentsCount, "%count%", strconv.Itoa(len(attachments)))
	return collapsibleSection(
		fmt.Sprintf(`<div class="section-icon" style="background:rgba(249,115,22,.1);color:#fb923c;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21.44 11.05l-9.19 9.19a6 6 0 01-8.49-8.49l9.19-9.19a4 4 0 015.66 5.66l-9.2 9.19a2 2 0 01-2.83-2.83l8.49-8.48"/></svg>
      </div>
      <span class="section-title">%s</span>
      <span class="section-count">%s</span>`, esc(labelAttachmentsTitle), esc(countLabel)),
		fmt.Sprintf(`<div class="attachments-grid">%s</div>
      <div class="attachments-note">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        %s
      </div>`, renderAttachmentCards(attachments), esc(labelAttachmentsNote)),
	)
}

func commentRoleClass(role string) string {
	switch strings.ToUpper(role) {
	case "COMPANY":
		return "comment-role-company"
	case "JUDGE":
		return "comment-role-judge"
	case "SUPPORT":
		return "comment-role-support"
	default:
		return "comment-role-bughunter"
	}
}

func avatarInitial(nickname string) string {
	trimmed := strings.TrimSpace(nickname)
	if trimmed == "" {
		return "?"
	}
	return strings.ToUpper(string([]rune(trimmed)[0]))
}

func renderComments(comments []models.Comment) string {
	if len(comments) == 0 {
		return ""
	}
	countLabel := strings.ReplaceAll(labelCommentsCount, "%count%", strconv.Itoa(len(comments)))
	var items strings.Builder
	for _, c := range comments {
		roleLabel := commentRoleLabelFA(c.By.Role)
		timeLabel := ""
		if c.CreatedAt != "" {
			timeLabel = fmt.Sprintf(`<span class="comment-time">%s</span>`, esc(formatDateTimeFA(c.CreatedAt)))
		}
		message := c.Message
		if c.RepliedTo != "" {
			message = "↩ " + c.RepliedTo + "\n\n" + message
		}
		nickname := strings.TrimSpace(firstNonEmpty(c.By.Nickname, c.By.DisplayName))
		if nickname == "" {
			nickname = roleLabel
		}
		fmt.Fprintf(&items, `<div class="comment-item">
        <div class="avatar">%s</div>
        <div class="comment-body">
          <div class="comment-meta">
            <span class="comment-author">%s</span>
            <span class="comment-role %s">%s</span>
            %s
          </div>
          <div class="comment-text md-content">%s</div>
        </div>
      </div>`, esc(avatarInitial(nickname)), esc(nickname), commentRoleClass(c.By.Role), esc(roleLabel), timeLabel, formatMarkdownHTML(message))
	}
	return collapsibleSection(
		fmt.Sprintf(`<div class="section-icon" style="background:rgba(148,163,184,.08);color:#94a3b8;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"/></svg>
      </div>
      <span class="section-title">%s</span>
      <span class="section-count">%s</span>`, esc(labelCommentsTitle), esc(countLabel)),
		items.String(),
	)
}

func renderAmendments(amendments []models.Amendment, outputDir string) string {
	if len(amendments) == 0 {
		return ""
	}
	countLabel := strings.ReplaceAll(labelAmendmentsCount, "%count%", strconv.Itoa(len(amendments)))
	var items strings.Builder
	for i, a := range amendments {
		content := a.Description
		if content == "" {
			content = a.Content
		}
		timeLabel := ""
		if !a.CreatedAt.IsZero() {
			timeLabel = fmt.Sprintf(`<span class="amendment-time" title="%s">%s</span>`, esc(labelAmendmentDate), esc(formatDateTimeFA(a.CreatedAt.Format(time.RFC3339))))
		} else if a.SubmissionDate != "" {
			timeLabel = fmt.Sprintf(`<span class="amendment-time" title="%s">%s</span>`, esc(labelAmendmentDate), esc(formatDateTimeFA(a.SubmissionDate)))
		}
		attachmentHTML := ""
		folderAtts := listAmendmentDirAttachments(outputDir, i+1)
		if len(folderAtts) > 0 {
			countLabelAtt := strings.ReplaceAll(labelAttachmentsCount, "%count%", strconv.Itoa(len(folderAtts)))
			attachmentHTML = fmt.Sprintf(`<div class="amendment-attachments">
                <div class="amendment-attachments-label">%s · %s</div>
                %s
              </div>`, esc(labelAmendmentAttachments), esc(countLabelAtt), renderAttachmentCards(folderAtts))
		}
		stageHTML := ""
		fromLabel := statusLabelFA(a.FromStage)
		toLabel := statusLabelFA(a.ToStage)
		if a.FromStage != "" && a.ToStage != "" {
			transition := fmt.Sprintf("%s: %s → %s: %s", labelAmendmentFrom, fromLabel, labelAmendmentTo, toLabel)
			stageHTML = fmt.Sprintf(`<div class="amendment-stage">
                <span class="amendment-stage-label">%s</span>
                %s
              </div>`, esc(labelAmendmentStage), esc(transition))
		}
		fmt.Fprintf(&items, `<div class="amendment-item">
        <div class="amendment-header">
          <span class="amendment-number">%s %d</span>
          %s
        </div>
        %s
        <div class="amendment-explanation-label">%s</div>
        <div class="amendment-explanation md-content">%s</div>
        %s
      </div>`, esc(labelAmendmentNumber), i+1, timeLabel, stageHTML, esc(labelAmendmentExplanation), formatMarkdownHTML(content), attachmentHTML)
	}
	return collapsibleSection(
		fmt.Sprintf(`<div class="section-icon" style="background:rgba(168,85,247,.1);color:#c084fc;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
      </div>
      <span class="section-title">%s</span>
      <span class="section-count">%s</span>`, esc(labelAmendmentsTitle), esc(countLabel)),
		items.String(),
	)
}

func themeToggleScript() string {
	return `<script>
(function () {
  var STORAGE_KEY = 'ravro-report-theme';
  var labels = { light: ` + strconv.Quote(labelThemeLight) + `, dark: ` + strconv.Quote(labelThemeDark) + ` };
  var sunIcon = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>';
  var moonIcon = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/></svg>';
  var btn = document.getElementById('theme-toggle');
  if (!btn) return;

  function getTheme() {
    var stored = localStorage.getItem(STORAGE_KEY);
    if (stored === 'dark') return 'dark';
    if (stored === 'light') return 'light';
    return 'light';
  }

  function updateButton(theme) {
    var isLight = theme === 'light';
    btn.setAttribute('aria-label', isLight ? labels.dark : labels.light);
    btn.innerHTML = (isLight ? moonIcon : sunIcon) + (isLight ? labels.dark : labels.light);
  }

  function applyTheme(theme) {
    if (theme === 'light') {
      document.documentElement.setAttribute('data-theme', 'light');
    } else {
      document.documentElement.removeAttribute('data-theme');
    }
    updateButton(theme);
  }

  btn.addEventListener('click', function () {
    var next = getTheme() === 'light' ? 'dark' : 'light';
    localStorage.setItem(STORAGE_KEY, next);
    applyTheme(next);
  });

  applyTheme(getTheme());
})();
</script>`
}
