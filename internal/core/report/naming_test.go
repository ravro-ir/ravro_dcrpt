package report

import (
	"testing"

	"ravro_dcrpt/pkg/models"
)

func TestOutputBundleNameMatchesHTMLStem(t *testing.T) {
	got := OutputBundleName(&models.Report{
		Slug:            "r2026-09-06-0002",
		CompanyUsername: "basalam",
		HunterUsername:  "mahyarm02",
	}, "ef9b30ee-96dc-4b25-9774-d760234ae4d6-ravro")
	want := "basalam__r2026-09-06-0002__mahyarm02"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestOutputBundleNameFallsBackAndSanitizes(t *testing.T) {
	got := OutputBundleName(&models.Report{
		CompanyUsername: `acme/corp`,
		HunterUsername:  "",
	}, "id:1")
	want := "acme_corp__id_1__hunter"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
