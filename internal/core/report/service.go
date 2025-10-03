package report

import (
	"fmt"
	"path/filepath"

	"ravro_dcrpt/internal/core/decrypt"
	"ravro_dcrpt/internal/ports"
	"ravro_dcrpt/pkg/models"
)

// Service handles report processing operations
type Service struct {
	decryptService *decrypt.Service
	pdfGenerator   ports.PDFGenerator
	storage        ports.StorageService
}

// NewService creates a new report service
func NewService(
	decryptService *decrypt.Service,
	pdfGenerator ports.PDFGenerator,
	storage ports.StorageService,
) *Service {
	return &Service{
		decryptService: decryptService,
		pdfGenerator:   pdfGenerator,
		storage:        storage,
	}
}

// ProcessReport processes a single report: decrypt and generate PDF
func (s *Service) ProcessReport(reportPath string, keyPath string, outputDir string) (*ProcessResult, error) {
	result := &ProcessResult{
		ReportPath: reportPath,
		ReportID:   s.decryptService.GetReportID(reportPath),
	}

	// Decrypt report
	report, err := s.decryptService.DecryptReport(reportPath, keyPath)
	if err != nil {
		result.Error = err
		return result, err
	}
	result.Report = report

	// Decrypt judgment
	judgment, err := s.decryptService.DecryptJudgment(reportPath, keyPath)
	if err != nil {
		// Judgment might not exist, don't fail
		judgment = nil
	}
	result.Judgment = judgment

	// Decrypt amendments
	amendments, err := s.decryptService.DecryptAmendment(reportPath, keyPath)
	if err != nil {
		// Amendments might not exist, don't fail
		amendments = nil
	}
	result.Amendments = amendments

	// Decrypt attachments (images, files, etc.)
	// Attachments should go to the same directory as the PDF
	pdfDir := filepath.Join(outputDir, report.Slug)
	if err := s.storage.CreateDir(pdfDir); err == nil {
		if err := s.decryptService.DecryptAttachments(reportPath, keyPath, pdfDir); err != nil {
			// Attachments decryption failure shouldn't stop PDF generation
			// Just log it
		}
	}

	// Save decrypted JSON files for debugging
	s.decryptService.SaveDecryptedJSON(reportPath, keyPath, outputDir)

	// Generate PDF
	pdfPath := s.generatePDFPath(outputDir, report)
	if err := s.pdfGenerator.GenerateReport(report, judgment, amendments, pdfPath); err != nil {
		result.Error = err
		return result, fmt.Errorf("failed to generate PDF: %w", err)
	}
	result.PDFPath = pdfPath
	result.Success = true

	return result, nil
}

// ProcessReports processes multiple reports
func (s *Service) ProcessReports(inputDir string, keyPath string, outputDir string) ([]*ProcessResult, error) {
	var reportPaths []string
	var results []*ProcessResult

	// Check if inputDir itself is a report directory (contains report/data.ravro)
	reportDataPath := filepath.Join(inputDir, "report", "data.ravro")
	if s.storage.FileExists(reportDataPath) {
		// This is a report directory itself
		reportPaths = append(reportPaths, inputDir)
	} else {
		// Find all zip files
		zipFiles, err := s.storage.ListFiles(inputDir, "*.zip")
		if err != nil {
			return nil, fmt.Errorf("failed to list zip files: %w", err)
		}

		// Extract zip files
		for _, zipFile := range zipFiles {
			reportID := s.decryptService.GetReportID(zipFile)
			extractPath := filepath.Join(inputDir, reportID)

			// Check if already extracted
			if !s.storage.FileExists(extractPath) {
				if err := s.decryptService.ProcessZipFile(zipFile, extractPath); err != nil {
					// Create a failed result for this ZIP file
					result := &ProcessResult{
						ReportPath: zipFile,
						ReportID:   reportID,
						Success:    false,
						Error:      fmt.Errorf("failed to extract ZIP: %w", err),
					}
					results = append(results, result)
					continue // Skip files that can't be extracted
				}
			}

			reportPaths = append(reportPaths, extractPath)
		}

		// Find existing ravro directories (already extracted)
		ravroFiles, err := s.storage.ListFiles(inputDir, "*.ravro")
		if err == nil {
			for _, ravroFile := range ravroFiles {
				reportDir := filepath.Dir(filepath.Dir(ravroFile)) // Go up two levels from data.ravro
				if !contains(reportPaths, reportDir) {
					reportPaths = append(reportPaths, reportDir)
				}
			}
		}
	}

	// Process each report
	for _, reportPath := range reportPaths {
		result, _ := s.ProcessReport(reportPath, keyPath, outputDir)
		results = append(results, result)
	}

	return results, nil
}

// generatePDFPath generates the output PDF file path
func (s *Service) generatePDFPath(outputDir string, report *models.Report) string {
	filename := fmt.Sprintf("%s__%s__%s.pdf",
		report.CompanyUsername,
		report.Slug,
		report.HunterUsername,
	)
	return filepath.Join(outputDir, report.Slug, filename)
}

// ProcessResult contains the result of processing a single report
type ProcessResult struct {
	ReportPath string
	ReportID   string
	Report     *models.Report
	Judgment   *models.Judgment
	Amendments []models.Amendment
	PDFPath    string
	Success    bool
	Error      error
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
