package report

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"ravro_dcrpt/internal/core/decrypt"
	"ravro_dcrpt/internal/ports"
	"ravro_dcrpt/pkg/models"
)

// ProgressCallback is called during report processing
type ProgressCallback func(current, total int, reportID string)

// Service handles report processing operations
type Service struct {
	decryptService *decrypt.Service
	htmlGenerator  ports.HTMLGenerator
	storage        ports.StorageService
}

// NewService creates a new report service
func NewService(
	decryptService *decrypt.Service,
	htmlGenerator ports.HTMLGenerator,
	storage ports.StorageService,
) *Service {
	return &Service{
		decryptService: decryptService,
		htmlGenerator:  htmlGenerator,
		storage:        storage,
	}
}

// ProcessReport processes a single report: decrypt and generate HTML
func (s *Service) ProcessReport(reportPath string, keyPath string, outputDir string) (*ProcessResult, error) {
	result := &ProcessResult{
		ReportPath: reportPath,
		ReportID:   s.decryptService.GetReportID(reportPath),
	}

	report, err := s.decryptService.DecryptReport(reportPath, keyPath)
	if err != nil {
		result.Error = err
		return result, err
	}
	result.Report = report

	judgment, err := s.decryptService.DecryptJudgment(reportPath, keyPath)
	if err != nil {
		judgment = nil
	}
	result.Judgment = judgment

	amendments, err := s.decryptService.DecryptAmendment(reportPath, keyPath)
	if err != nil {
		amendments = nil
	}
	result.Amendments = amendments

	bundle := OutputBundleName(report, result.ReportID)
	result.ReportID = bundle
	htmlDir := filepath.Join(outputDir, bundle)
	if err := s.storage.CreateDir(htmlDir); err == nil {
		_ = s.decryptService.DecryptAttachments(reportPath, keyPath, htmlDir)
		_ = s.decryptService.CopyReportSidecars(reportPath, htmlDir)
	}

	htmlPath := filepath.Join(htmlDir, bundle+".html")
	if err := s.htmlGenerator.GenerateReport(report, judgment, amendments, htmlPath); err != nil {
		result.Error = err
		return result, fmt.Errorf("failed to generate HTML: %w", err)
	}
	result.HTMLPath = htmlPath
	result.Success = true

	return result, nil
}

// OutputBundleName is company__slug__hunter — same stem as the HTML file and output folder.
func OutputBundleName(report *models.Report, fallbackID string) string {
	id := sanitizeNamePart(firstNonEmpty(report.Slug, fallbackID, "report"))
	company := sanitizeNamePart(firstNonEmpty(report.CompanyUsername, "company"))
	hunter := sanitizeNamePart(firstNonEmpty(report.HunterUsername, "hunter"))
	return company + "__" + id + "__" + hunter
}

func sanitizeNamePart(value string) string {
	value = strings.TrimSpace(value)
	var b strings.Builder
	b.Grow(len(value))
	for _, r := range value {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			b.WriteByte('_')
		default:
			if unicode.IsControl(r) {
				b.WriteByte('_')
				continue
			}
			b.WriteRune(r)
		}
	}
	out := strings.Trim(b.String(), ". ")
	if out == "" {
		return "unknown"
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// ProcessReports processes multiple reports
func (s *Service) ProcessReports(inputDir string, keyPath string, outputDir string) ([]*ProcessResult, error) {
	return s.ProcessReportsWithProgress(inputDir, keyPath, outputDir, nil)
}

// ProcessReportsWithProgress processes multiple reports with progress callback
func (s *Service) ProcessReportsWithProgress(inputDir string, keyPath string, outputDir string, progressCallback ProgressCallback) ([]*ProcessResult, error) {
	var reportPaths []string
	var results []*ProcessResult

	reportDataPath := filepath.Join(inputDir, "report", "data.ravro")
	if s.storage.FileExists(reportDataPath) {
		reportPaths = append(reportPaths, inputDir)
	} else {
		zipFiles, err := s.storage.ListFiles(inputDir, "*.zip")
		if err != nil {
			return nil, fmt.Errorf("failed to list zip files: %w", err)
		}

		for _, zipFile := range zipFiles {
			reportID := s.decryptService.GetReportID(zipFile)
			extractPath := filepath.Join(inputDir, reportID)

			if !s.storage.FileExists(extractPath) {
				if err := s.decryptService.ProcessZipFile(zipFile, extractPath); err != nil {
					result := &ProcessResult{
						ReportPath: zipFile,
						ReportID:   reportID,
						Success:    false,
						Error:      fmt.Errorf("failed to extract ZIP: %w", err),
					}
					results = append(results, result)
					continue
				}
			}

			reportPaths = append(reportPaths, extractPath)
		}

		if len(reportPaths) == 0 {
			ravroFiles, err := s.storage.ListFiles(inputDir, "*.ravro")
			if err == nil {
				for _, ravroFile := range ravroFiles {
					reportDir := filepath.Dir(filepath.Dir(ravroFile))
					if !contains(reportPaths, reportDir) {
						reportPaths = append(reportPaths, reportDir)
					}
				}
			}
		}
	}

	total := len(reportPaths)

	for i, reportPath := range reportPaths {
		reportID := s.decryptService.GetReportID(reportPath)

		if progressCallback != nil {
			progressCallback(i+1, total, reportID)
		}

		result, _ := s.ProcessReport(reportPath, keyPath, outputDir)
		results = append(results, result)
	}

	return results, nil
}

// ProcessResult contains the result of processing a single report
type ProcessResult struct {
	ReportPath string
	ReportID   string
	Report     *models.Report
	Judgment   *models.Judgment
	Amendments []models.Amendment
	HTMLPath   string
	Success    bool
	Error      error
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
