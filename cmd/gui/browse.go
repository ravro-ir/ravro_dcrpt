package main

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/dialog"
	fynestorage "fyne.io/fyne/v2/storage"
)

func binaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		wd, _ := os.Getwd()
		return wd
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return filepath.Dir(exe)
}

func resolveBrowseDir(preferred, base string) string {
	tryDir := func(p string) string {
		info, err := os.Stat(p)
		if err != nil {
			return ""
		}
		if info.IsDir() {
			return p
		}
		return filepath.Dir(p)
	}
	if p := strings.TrimSpace(preferred); p != "" {
		if !filepath.IsAbs(p) {
			p = filepath.Join(base, p)
		}
		if dir := tryDir(p); dir != "" {
			return dir
		}
	}
	if dir := tryDir(base); dir != "" {
		return dir
	}
	wd, _ := os.Getwd()
	return wd
}

func setBrowseLocation(fd *dialog.FileDialog, preferred string) {
	dir := resolveBrowseDir(preferred, binaryDir())
	list, err := fynestorage.ListerForURI(fynestorage.NewFileURI(dir))
	if err != nil {
		return
	}
	fd.SetLocation(list)
}
