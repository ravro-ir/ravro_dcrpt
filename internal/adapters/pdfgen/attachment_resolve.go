package pdfgen

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"ravro_dcrpt/pkg/models"
)

var numberedAttachmentName = regexp.MustCompile(`^(\d+)\.\s*(.+)$`)

type checksumEntry struct {
	Index int
	Name  string
}

type resolvedAttachment struct {
	relPath string
	display string
}

func isMetaAttachmentFile(name string) bool {
	lower := strings.ToLower(name)
	switch lower {
	case "sha1_checksum.txt", "data.json", "data.com", "data.ravro", "data.com.ravro", "comments.json", "report_info.json":
		return true
	}
	if strings.HasSuffix(lower, ".html") || strings.HasSuffix(lower, ".htm") {
		return true
	}
	return false
}

func parseSHA1Checksum(content string) []checksumEntry {
	var entries []checksumEntry
	var current checksumEntry
	haveIndex := false

	flush := func() {
		if haveIndex && strings.TrimSpace(current.Name) != "" {
			entries = append(entries, current)
		}
		current = checksumEntry{}
		haveIndex = false
	}

	sc := bufio.NewScanner(strings.NewReader(content))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			flush()
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch strings.ToLower(key) {
		case "index":
			flush()
			n, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			current.Index = n
			haveIndex = true
		case "name":
			current.Name = value
		}
	}
	flush()
	return entries
}

func loadSHA1Checksum(dirs ...string) []checksumEntry {
	var candidates []string
	for _, dir := range dirs {
		if strings.TrimSpace(dir) == "" {
			continue
		}
		candidates = append(candidates,
			filepath.Join(dir, "sha1_checksum.txt"),
			filepath.Join(dir, "report", "sha1_checksum.txt"),
		)
	}
	for _, p := range candidates {
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if entries := parseSHA1Checksum(string(raw)); len(entries) > 0 {
			return entries
		}
	}
	return nil
}

func stripAttachmentNumberPrefix(name string) string {
	m := numberedAttachmentName.FindStringSubmatch(name)
	if len(m) == 3 {
		return m[2]
	}
	return name
}

func paddedAttachmentIndex(index int) string {
	return fmt.Sprintf("%03d", index)
}

func checksumRelPath(entry checksumEntry) string {
	return path.Join("report", paddedAttachmentIndex(entry.Index)+". "+entry.Name)
}

func isRemoteURL(raw string) bool {
	lower := strings.ToLower(strings.TrimSpace(raw))
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "javascript:") ||
		strings.HasPrefix(lower, "data:")
}

// encodeRelPath keeps local filenames as they appear on disk.
// Percent-encoding spaces breaks file:// resolution in Chrome on Linux.
func encodeRelPath(rel string) string {
	rel = path.Clean(strings.ReplaceAll(rel, "\\", "/"))
	if rel == "." {
		return ""
	}
	return rel
}

func listDirAttachments(dir, relPrefix string) []resolvedAttachment {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []resolvedAttachment
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if isMetaAttachmentFile(name) {
			continue
		}
		out = append(out, resolvedAttachment{
			relPath: path.Join(relPrefix, name),
			display: name,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].relPath < out[j].relPath
	})
	return out
}

func listReportDirAttachments(outputDir string) []resolvedAttachment {
	return listDirAttachments(filepath.Join(outputDir, "report"), "report")
}

func listAmendmentDirAttachments(outputDir string, index int) []resolvedAttachment {
	var out []resolvedAttachment
	for _, spelling := range []string{"amendment", "ammendment"} {
		rel := fmt.Sprintf("%s-%d", spelling, index)
		found := listDirAttachments(filepath.Join(outputDir, rel), rel)
		out = append(out, found...)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].relPath < out[j].relPath
	})
	return out
}

func findDiskAttachment(disk []resolvedAttachment, filename string) *resolvedAttachment {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil
	}
	base := path.Base(filename)
	plain := stripAttachmentNumberPrefix(base)
	for i := range disk {
		item := disk[i]
		itemBase := path.Base(item.relPath)
		if itemBase == filename || itemBase == base || item.relPath == filename {
			return &item
		}
		if stripAttachmentNumberPrefix(itemBase) == plain || stripAttachmentNumberPrefix(itemBase) == base {
			return &item
		}
	}
	return nil
}

func findChecksumEntry(entries []checksumEntry, filename string) *checksumEntry {
	filename = strings.TrimSpace(filename)
	plain := stripAttachmentNumberPrefix(path.Base(filename))
	for i := range entries {
		if entries[i].Name == filename || entries[i].Name == plain || entries[i].Name == path.Base(filename) {
			return &entries[i]
		}
	}
	return nil
}

func resolveJSONAttachment(att models.Attachment, disk []resolvedAttachment, checksum []checksumEntry) resolvedAttachment {
	name := strings.TrimSpace(att.Filename)
	if found := findDiskAttachment(disk, name); found != nil {
		return *found
	}
	url := strings.TrimSpace(att.URL)
	if url != "" && !isRemoteURL(url) {
		if found := findDiskAttachment(disk, url); found != nil {
			return *found
		}
	}
	if entry := findChecksumEntry(checksum, name); entry != nil {
		rel := checksumRelPath(*entry)
		return resolvedAttachment{relPath: rel, display: path.Base(rel)}
	}
	if url != "" && !isRemoteURL(url) && strings.Contains(url, "/") {
		return resolvedAttachment{relPath: strings.ReplaceAll(url, "\\", "/"), display: firstNonEmpty(name, path.Base(url))}
	}
	if name != "" && strings.Contains(name, "/") {
		return resolvedAttachment{relPath: name, display: path.Base(name)}
	}
	if name == "" {
		return resolvedAttachment{}
	}
	return resolvedAttachment{relPath: path.Join("report", name), display: name}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func resolveReportAttachments(outputDir string, jsonAtts []models.Attachment) []resolvedAttachment {
	disk := listReportDirAttachments(outputDir)
	if len(disk) > 0 {
		return disk
	}
	checksum := loadSHA1Checksum(outputDir)
	var out []resolvedAttachment
	seen := map[string]bool{}
	for _, att := range jsonAtts {
		resolved := resolveJSONAttachment(att, disk, checksum)
		if resolved.relPath == "" || seen[resolved.relPath] {
			continue
		}
		seen[resolved.relPath] = true
		out = append(out, resolved)
	}
	return out
}
