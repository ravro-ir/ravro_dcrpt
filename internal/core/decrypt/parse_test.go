package decrypt

import (
	"testing"

	"ravro_dcrpt/pkg/models"
)

func TestDetectSchema(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want int
	}{
		{
			name: "v1 flat has no schema version",
			raw:  `{"title":"x","hunterUsername":"h"}`,
			want: models.SchemaVersionV1,
		},
		{
			name: "v2 camelCase schemaVersion",
			raw:  `{"schemaVersion":2,"payload":{"title":"x"}}`,
			want: models.SchemaVersionV2,
		},
		{
			name: "v2 snake_case schema_version",
			raw:  `{"schema_version":2,"payload":{"explanation":"x"}}`,
			want: models.SchemaVersionV2,
		},
		{
			name: "payload alone is treated as v2",
			raw:  `{"payload":{"title":"x"}}`,
			want: models.SchemaVersionV2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectSchema([]byte(tt.raw)); got != tt.want {
				t.Fatalf("DetectSchema() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestParseReportJSON_V1(t *testing.T) {
	raw := []byte(`{
		"title": "XSS",
		"description": "desc",
		"scenario": "steps",
		"submissionDate": "2026-01-01 12:00:00",
		"hunterUsername": "alice",
		"companyUsername": "acme",
		"targetName": "Web App",
		"attachmentsCount": 0,
		"slug": "r2026-01-01-0001",
		"category": {"title": "XSS", "impact": "high"},
		"ips": "1.2.3.4",
		"cvss": "CVSS:3.1/AV:N",
		"urls": "https://example.com",
		"dateFrom": "2026-01-01",
		"dateTo": "2026-01-02"
	}`)

	report, version, err := ParseReportJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if version != models.SchemaVersionV1 {
		t.Fatalf("version = %d, want 1", version)
	}
	if report.Title != "XSS" || report.HunterUsername != "alice" || report.CompanyUsername != "acme" {
		t.Fatalf("unexpected identities: %+v", report)
	}
	if report.ReportInfo.Details.Target != "Web App" {
		t.Fatalf("target = %q, want Web App", report.ReportInfo.Details.Target)
	}
	if report.ReportInfo.Details.Cvss.Hunter.Vector != "CVSS:3.1/AV:N" {
		t.Fatalf("hunter cvss = %q", report.ReportInfo.Details.Cvss.Hunter.Vector)
	}
}

func TestParseReportJSON_V2(t *testing.T) {
	raw := []byte(`{
		"schemaVersion": 2,
		"createdAt": "2026-01-01T12:00:00Z",
		"payload": {
			"title": "SQLi",
			"description": "desc",
			"scenario": "repro",
			"submission_date": "2026-01-01 12:00:00",
			"slug": "r2026-01-01-0001",
			"attachments_count": 0,
			"ips": "1.2.3.4,5.6.7.8",
			"urls": "https://example.com",
			"cvss": "CVSS:3.1/AV:N/AC:L",
			"date_from": "2026-01-01",
			"date_to": "2026-01-02",
			"bughunter": {"id": "u1", "slug": "alice", "title": "Alice"},
			"company": {"id": "c1", "slug": "acme", "title": "Acme Co"},
			"workspace": {"id": "w1", "slug": "ws", "title": "WS"},
			"target": {"id": "t1", "slug": "web", "title": "Web App"}
		}
	}`)

	report, version, err := ParseReportJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if version != models.SchemaVersionV2 {
		t.Fatalf("version = %d, want 2", version)
	}
	if report.HunterUsername != "alice" {
		t.Fatalf("hunter = %q, want alice", report.HunterUsername)
	}
	if report.CompanyUsername != "acme" {
		t.Fatalf("company = %q, want acme", report.CompanyUsername)
	}
	if report.SubmissionDate != "2026-01-01 12:00:00" || report.DateFrom != "2026-01-01" {
		t.Fatalf("dates: %+v", report)
	}
	if report.ReportInfo.Details.Target != "Web App" {
		t.Fatalf("target = %q", report.ReportInfo.Details.Target)
	}
	if report.ReportInfo.Details.Cvss.Hunter.Vector != "CVSS:3.1/AV:N/AC:L" {
		t.Fatalf("cvss = %q", report.ReportInfo.Details.Cvss.Hunter.Vector)
	}
}

func TestParseAmendmentJSON(t *testing.T) {
	t.Run("v1 numeric reportId", func(t *testing.T) {
		raw := []byte(`{
			"reportId": 123,
			"attachmentsCount": 1,
			"description": "more details",
			"hunterUsername": "alice",
			"companyUsername": "acme",
			"submissionDate": "2026-01-01 12:00:00"
		}`)
		a, err := ParseAmendmentJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
		if a.ReportID != "123" {
			t.Fatalf("reportId = %q, want 123", a.ReportID)
		}
		if a.Description != "more details" {
			t.Fatalf("description = %q", a.Description)
		}
	})

	t.Run("v2 explanation and slug", func(t *testing.T) {
		raw := []byte(`{
			"schema_version": 2,
			"created_at": "2026-01-01T12:00:00Z",
			"payload": {
				"explanation": "follow-up",
				"attachments": [
					{"attachment_id": "uuid", "ext_hash": "EXTabc"}
				],
				"report_slug": "r2026-01-01-0001",
				"from_stage": "INITIAL_REJECT",
				"to_stage": "INITIAL_REVIEW"
			}
		}`)
		a, err := ParseAmendmentJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
		if a.Description != "follow-up" {
			t.Fatalf("description = %q, want follow-up", a.Description)
		}
		if a.ReportID != "r2026-01-01-0001" {
			t.Fatalf("reportId = %q", a.ReportID)
		}
		if a.AttachmentsCount != 1 {
			t.Fatalf("attachmentsCount = %d", a.AttachmentsCount)
		}
	})
}

func TestParseReportInfoJSON(t *testing.T) {
	t.Run("v1 hunter flat and camelCase tags", func(t *testing.T) {
		raw := []byte(`{
			"tags": [
				{
					"infoTitle": "XSS",
					"infoDescription": "desc",
					"infoSolution": "fix",
					"infoMore": "https://owasp.org"
				}
			],
			"details": {
				"currentStatus": "در حال بررسی",
				"target": "Web",
				"judges": [{"name": "j1"}],
				"attachments": [{"filename": "a.png", "type": "image/png"}],
				"cvss": {
					"hunter": {"vector": "CVSS:3.1/AV:N", "score": "7.5"},
					"judge": {"vector": "CVSS:3.1/AV:L", "score": "4.0"}
				}
			}
		}`)
		info, err := ParseReportInfoJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
		if len(info.Tags) != 1 || info.Tags[0].InfoTitle != "XSS" || info.Tags[0].InfoSolution != "fix" {
			t.Fatalf("tags: %+v", info.Tags)
		}
		if info.Details.Cvss.Hunter.Vector != "CVSS:3.1/AV:N" || info.Details.Cvss.Hunter.Score != "7.5" {
			t.Fatalf("hunter cvss: %+v", info.Details.Cvss.Hunter)
		}
		if info.Details.Cvss.Judge.Score != "4.0" {
			t.Fatalf("judge cvss: %+v", info.Details.Cvss.Judge)
		}
	})

	t.Run("v2 bughunter nested prefers cvss40", func(t *testing.T) {
		raw := []byte(`{
			"tags": [],
			"details": {
				"currentStatus": "INITIAL_REVIEW",
				"target": "Target title",
				"judges": [],
				"attachments": [{"filename": "safe-name.pdf", "type": "application/pdf"}],
				"cvss": {
					"bughunter": {
						"cvss31": {"vector": "CVSS:3.1/AV:N", "score": "5.0"},
						"cvss40": {"vector": "CVSS:4.0/AV:N", "score": "8.1"}
					},
					"judge": {
						"cvss31": {"vector": null, "score": null},
						"cvss40": {"vector": null, "score": null}
					}
				}
			}
		}`)
		info, err := ParseReportInfoJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
		if info.Details.Cvss.Hunter.Vector != "CVSS:4.0/AV:N" || info.Details.Cvss.Hunter.Score != "8.1" {
			t.Fatalf("expected cvss40, got %+v", info.Details.Cvss.Hunter)
		}
		if info.Details.Cvss.Judge.Vector != "" {
			t.Fatalf("judge should be empty on nulls, got %+v", info.Details.Cvss.Judge)
		}
	})

	t.Run("v2 falls back to cvss31 when 4.0 empty", func(t *testing.T) {
		raw := []byte(`{
			"details": {
				"cvss": {
					"bughunter": {
						"cvss31": {"vector": "CVSS:3.1/AV:A", "score": "6.1"},
						"cvss40": {"vector": null, "score": null}
					}
				}
			}
		}`)
		info, err := ParseReportInfoJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
		if info.Details.Cvss.Hunter.Vector != "CVSS:3.1/AV:A" {
			t.Fatalf("got %+v", info.Details.Cvss.Hunter)
		}
	})

	t.Run("schema v2 export block from backend-api", func(t *testing.T) {
		raw := []byte(`{
			"schema_version": 2,
			"tags": [],
			"details": {
				"currentStatus": "ACCEPTED",
				"target": "Web",
				"cvss": {
					"company": {"cvss31": {"vector": "CVSS:3.1/AV:N", "score": "4.0"}},
					"final": {"cvss40": {"vector": "CVSS:4.0/AV:N", "score": "6.8"}}
				}
			},
			"export": {
				"updated_at": "2026-10-01T10:00:00Z",
				"bughunter": {"title": "hunter-one", "slug": "h1", "url": "https://ravro.ir/fa/hunters/h1"},
				"company": {"title": "Acme", "slug": "acme", "url": "https://ravro.ir/fa/companies/acme"},
				"target": {"title": "Shop", "slug": "shop", "version": 3, "rules": "## Rule"},
				"bounty": {"judge": "1000000", "company": "1200000", "final": "1100000", "paid_at": "2026-10-02T12:00:00Z"},
				"duplicate_of_report_slug": "r2026-10-01-0004",
				"timeline": [{"at": "2026-10-01T11:00:00Z", "from_stage": "INITIAL", "to_stage": "JUDGEMENT"}],
				"appeal": {"status": "PENDING", "target_role": "COMPANY", "message": "please review"}
			}
		}`)
		info, err := ParseReportInfoJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
		if !info.Export.Present || info.Export.Company.Title != "Acme" {
			t.Fatalf("export: %+v", info.Export)
		}
		if info.Export.Target.Version != "3" || info.Export.Bounty.Judge != "1000000" {
			t.Fatalf("target/bounty: %+v %+v", info.Export.Target, info.Export.Bounty)
		}
		if info.Details.Cvss.Company.Score != "4.0" || info.Details.Cvss.Final.Score != "6.8" {
			t.Fatalf("cvss company/final: %+v", info.Details.Cvss)
		}
		if len(info.Export.Timeline) != 1 || info.Export.Appeal == nil || info.Export.Appeal.Status != "PENDING" {
			t.Fatalf("timeline/appeal: %+v", info.Export)
		}
	})
}

func TestMergeReportInfoKeepsPayloadCVSS(t *testing.T) {
	report := &models.Report{
		Cvss: "CVSS:3.1/FROM-PAYLOAD",
		ReportInfo: models.InfoReport{
			Details: models.ReportDetails{
				Target: "seed-target",
				Cvss:   models.CVSSInfo{Hunter: models.CVSSScore{Vector: "CVSS:3.1/FROM-PAYLOAD"}},
			},
		},
	}
	info := &models.InfoReport{
		Details: models.ReportDetails{
			CurrentStatus: "OPEN",
			Cvss:          models.CVSSInfo{Hunter: models.CVSSScore{}},
		},
	}
	MergeReportInfo(report, info)
	if report.ReportInfo.Details.Target != "seed-target" {
		t.Fatalf("target lost: %q", report.ReportInfo.Details.Target)
	}
	if report.ReportInfo.Details.Cvss.Hunter.Vector != "CVSS:3.1/FROM-PAYLOAD" {
		t.Fatalf("cvss lost: %+v", report.ReportInfo.Details.Cvss.Hunter)
	}
	if report.ReportInfo.Details.CurrentStatus != "OPEN" {
		t.Fatalf("status = %q", report.ReportInfo.Details.CurrentStatus)
	}
}

func TestPathClassifiers(t *testing.T) {
	tests := []struct {
		path       string
		report     bool
		judgment   bool
		amendment  bool
		structured bool
	}{
		{"encrypt/r1/report/data.ravro", true, false, false, true},
		{"encrypt/r1/judgment/data.ravro", false, true, false, true},
		{"encrypt/r1/amendment-1/data.ravro", false, false, true, true},
		{"encrypt/r1/ammendment-1/data.ravro", false, false, true, true},
		{"encrypt/r1/ammendment-1/data.com.ravro", false, false, true, true},
		{`encrypt\r1\ammendment-1\data.com.ravro`, false, false, true, true},
		{"encrypt/r1/report/001. file.pdf.ravro", false, false, false, false},
		{"encrypt/r1/amendment-1/002. extra.pdf.ravro", false, false, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := IsReportPayloadFile(tt.path); got != tt.report {
				t.Errorf("IsReportPayloadFile = %v, want %v", got, tt.report)
			}
			if got := IsJudgmentPayloadFile(tt.path); got != tt.judgment {
				t.Errorf("IsJudgmentPayloadFile = %v, want %v", got, tt.judgment)
			}
			if got := IsAmendmentPayloadFile(tt.path); got != tt.amendment {
				t.Errorf("IsAmendmentPayloadFile = %v, want %v", got, tt.amendment)
			}
			if got := IsStructuredDataFile(tt.path); got != tt.structured {
				t.Errorf("IsStructuredDataFile = %v, want %v", got, tt.structured)
			}
		})
	}
}

func TestParseCommentsJSON(t *testing.T) {
	raw := []byte(`{
		"schema_version": 1,
		"comments": [
			{"message": "hello", "by": {"nickname": "alice", "role": "BUGHUNTER"}, "created_at": "2026-01-01T12:00:00Z"},
			{"message": "reply", "replied_to": "hello", "by": {"nickname": "bob", "role": "COMPANY"}, "created_at": "2026-01-01T13:00:00Z"}
		]
	}`)
	comments, err := ParseCommentsJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 || comments[1].RepliedTo != "hello" {
		t.Fatalf("%+v", comments)
	}
}

func TestParseCommentsJSONDisplayName(t *testing.T) {
	raw := []byte(`{
		"schema_version": 2,
		"comments": [
			{"message": "note", "by": {"role": "SUPPORT", "display_name": "wwrbvettony4"}, "created_at": "2026-10-01T08:00:00Z"},
			{"message": "hunter", "by": {"role": "BUGHUNTER", "display_name": "hunter-one"}, "created_at": "2026-10-01T09:00:00Z"}
		]
	}`)
	comments, err := ParseCommentsJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 {
		t.Fatalf("len=%d", len(comments))
	}
	if comments[0].By.Nickname != "wwrbvettony4" || comments[0].By.Role != "SUPPORT" {
		t.Fatalf("support: %+v", comments[0].By)
	}
	if comments[1].By.Nickname != "hunter-one" {
		t.Fatalf("hunter: %+v", comments[1].By)
	}
}
