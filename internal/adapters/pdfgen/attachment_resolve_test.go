package pdfgen

import (
	"os"
	"path/filepath"
	"testing"

	"ravro_dcrpt/pkg/models"
)

func TestParseSHA1Checksum(t *testing.T) {
	raw := `Index: 1
Name: files.zip
CreatedAt: 2026-01-01 12:00:00
Type: application/zip
SHA1 Checksum: abc

Index: 2
Name: Screenshot_217.png
Type: image/png
SHA1 Checksum: def
`
	entries := parseSHA1Checksum(raw)
	if len(entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(entries))
	}
	if entries[0].Index != 1 || entries[0].Name != "files.zip" {
		t.Fatalf("first entry = %+v", entries[0])
	}
	if entries[1].Index != 2 || entries[1].Name != "Screenshot_217.png" {
		t.Fatalf("second entry = %+v", entries[1])
	}
}

func TestResolveReportAttachmentsUsesChecksumPrefix(t *testing.T) {
	dir := t.TempDir()
	checksum := "Index: 1\nName: files.zip\n\nIndex: 2\nName: Screenshot_217.png\n"
	if err := os.WriteFile(filepath.Join(dir, "sha1_checksum.txt"), []byte(checksum), 0644); err != nil {
		t.Fatal(err)
	}

	got := resolveReportAttachments(dir, []models.Attachment{
		{Filename: "files.zip"},
		{Filename: "Screenshot_217.png"},
	})
	if len(got) != 2 {
		t.Fatalf("got %d attachments: %+v", len(got), got)
	}
	if got[0].relPath != "report/001. files.zip" {
		t.Fatalf("first path = %q", got[0].relPath)
	}
	if got[1].relPath != "report/002. Screenshot_217.png" {
		t.Fatalf("second path = %q", got[1].relPath)
	}
}

func TestResolveReportAttachmentsPrefersOnDiskNumberedFiles(t *testing.T) {
	dir := t.TempDir()
	reportDir := filepath.Join(dir, "report")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reportDir, "001. files.zip"), []byte("zip"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reportDir, "002. Screenshot_217.png"), []byte("png"), 0644); err != nil {
		t.Fatal(err)
	}

	got := resolveReportAttachments(dir, []models.Attachment{
		{Filename: "files.zip", URL: "https://cdn.example/files.zip"},
	})
	if len(got) != 2 {
		t.Fatalf("got %+v", got)
	}
	if got[0].relPath != "report/001. files.zip" || got[1].relPath != "report/002. Screenshot_217.png" {
		t.Fatalf("got %+v", got)
	}
	if encodeRelPath(got[0].relPath) != "report/001. files.zip" {
		t.Fatalf("encoded = %q", encodeRelPath(got[0].relPath))
	}
}
