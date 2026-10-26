package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveBrowseDirPrefersExistingRelativeFolder(t *testing.T) {
	base := t.TempDir()
	encrypt := filepath.Join(base, "encrypt")
	if err := os.Mkdir(encrypt, 0755); err != nil {
		t.Fatal(err)
	}
	got := resolveBrowseDir("encrypt", base)
	if got != encrypt {
		t.Fatalf("got %q, want %q", got, encrypt)
	}
}

func TestResolveBrowseDirFallsBackToBinaryDir(t *testing.T) {
	base := t.TempDir()
	got := resolveBrowseDir("missing", base)
	if got != base {
		t.Fatalf("got %q, want %q", got, base)
	}
}

func TestResolveBrowseDirUsesParentOfFile(t *testing.T) {
	base := t.TempDir()
	key := filepath.Join(base, "company.pem")
	if err := os.WriteFile(key, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	got := resolveBrowseDir(key, base)
	if got != base {
		t.Fatalf("got %q, want %q", got, base)
	}
}
