package decrypt

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"ravro_dcrpt/internal/ports"
	"ravro_dcrpt/pkg/models"
)

// Service handles decryption operations
type Service struct {
	crypto  ports.CryptoService
	storage ports.StorageService
}

// NewService creates a new decryption service
func NewService(crypto ports.CryptoService, storage ports.StorageService) *Service {
	return &Service{
		crypto:  crypto,
		storage: storage,
	}
}

// DecryptReport decrypts a report from a ravro file
func (s *Service) DecryptReport(reportPath string, privateKeyPath string) (*models.Report, error) {
	// Find report data file
	ravroFiles, err := s.storage.ListFiles(reportPath, "*.ravro")
	if err != nil {
		return nil, fmt.Errorf("failed to list ravro files: %w", err)
	}

	var reportFile string
	for _, file := range ravroFiles {
		if IsReportPayloadFile(file) {
			reportFile = file
			break
		}
	}

	if reportFile == "" {
		return nil, fmt.Errorf("report data file not found")
	}

	// Decrypt the report file
	decrypted, err := s.decryptFile(reportFile, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt report: %w", err)
	}

	report, _, err := ParseReportJSON(decrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to parse report JSON: %w", err)
	}

	infoFiles, _ := s.storage.ListFiles(reportPath, "report_info.json")
	if len(infoFiles) > 0 {
		infoData, err := s.storage.ReadFile(infoFiles[0])
		if err == nil {
			if info, err := ParseReportInfoJSON(infoData); err == nil {
				MergeReportInfo(report, info)
			}
		}
	}

	commentFiles, _ := s.storage.ListFiles(reportPath, "comments.json")
	if len(commentFiles) > 0 {
		commentData, err := s.storage.ReadFile(commentFiles[0])
		if err == nil {
			if comments, err := ParseCommentsJSON(commentData); err == nil {
				report.Comments = comments
			}
		}
	}

	return report, nil
}

// SaveDecryptedJSON saves decrypted JSON for debugging
func (s *Service) SaveDecryptedJSON(reportPath string, privateKeyPath string, outputDir string) error {
	// Get report ID
	reportID := s.GetReportID(reportPath)
	debugDir := filepath.Join(outputDir, reportID, "debug")

	// Create debug directory
	if err := s.storage.CreateDir(debugDir); err != nil {
		return err
	}

	// Decrypt and save report JSON
	if report, err := s.DecryptReport(reportPath, privateKeyPath); err == nil {
		reportJSON, _ := json.MarshalIndent(report, "", "  ")
		s.storage.WriteFile(filepath.Join(debugDir, "report.json"), reportJSON)
	}

	// Decrypt and save judgment JSON
	if judgment, err := s.DecryptJudgment(reportPath, privateKeyPath); err == nil {
		judgmentJSON, _ := json.MarshalIndent(judgment, "", "  ")
		s.storage.WriteFile(filepath.Join(debugDir, "judgment.json"), judgmentJSON)
	}

	// Decrypt and save amendments JSON
	if amendments, err := s.DecryptAmendment(reportPath, privateKeyPath); err == nil && len(amendments) > 0 {
		amendmentsJSON, _ := json.MarshalIndent(amendments, "", "  ")
		s.storage.WriteFile(filepath.Join(debugDir, "amendments.json"), amendmentsJSON)
	}

	return nil
}

// DecryptJudgment decrypts judgment data
func (s *Service) DecryptJudgment(reportPath string, privateKeyPath string) (*models.Judgment, error) {
	// Find judgment data file
	ravroFiles, err := s.storage.ListFiles(reportPath, "*.ravro")
	if err != nil {
		return nil, fmt.Errorf("failed to list ravro files: %w", err)
	}

	var judgmentFile string
	for _, file := range ravroFiles {
		if IsJudgmentPayloadFile(file) {
			judgmentFile = file
			break
		}
	}

	if judgmentFile == "" {
		return nil, fmt.Errorf("judgment data file not found")
	}

	// Decrypt the judgment file
	decrypted, err := s.decryptFile(judgmentFile, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt judgment: %w", err)
	}

	// Parse JSON
	var judgment models.Judgment
	if err := json.Unmarshal(decrypted, &judgment); err != nil {
		return nil, fmt.Errorf("failed to parse judgment JSON: %w", err)
	}

	return &judgment, nil
}

// DecryptAmendment decrypts amendment data
func (s *Service) DecryptAmendment(reportPath string, privateKeyPath string) ([]models.Amendment, error) {
	// Find amendment data files
	ravroFiles, err := s.storage.ListFiles(reportPath, "*.ravro")
	if err != nil {
		return nil, fmt.Errorf("failed to list ravro files: %w", err)
	}

	var amendments []models.Amendment

	for _, file := range ravroFiles {
		if !IsAmendmentPayloadFile(file) {
			continue
		}
		decrypted, err := s.decryptFile(file, privateKeyPath)
		if err != nil {
			continue
		}
		amendment, err := ParseAmendmentJSON(decrypted)
		if err != nil {
			continue
		}
		amendments = append(amendments, *amendment)
	}

	return amendments, nil
}

// decryptFile decrypts a single file
func (s *Service) decryptFile(filePath string, privateKeyPath string) ([]byte, error) {
	// Read encrypted file
	encrypted, err := s.storage.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Decrypt
	decrypted, err := s.crypto.DecryptPKCS7(encrypted, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return decrypted, nil
}

// ProcessZipFile extracts and processes a zip file
func (s *Service) ProcessZipFile(zipPath string, extractPath string) error {
	// Extract zip file
	if err := s.storage.ExtractZip(zipPath, extractPath); err != nil {
		return fmt.Errorf("failed to extract zip file: %w", err)
	}

	return nil
}

// ValidateKey validates a private key file
func (s *Service) ValidateKey(privateKeyPath string) error {
	return s.crypto.ValidatePrivateKey(privateKeyPath)
}

// DecryptAttachments decrypts all attachment files (images, etc.)
func (s *Service) DecryptAttachments(reportPath string, privateKeyPath string, outputPath string) error {
	// Find all ravro files
	ravroFiles, err := s.storage.ListFiles(reportPath, "*.ravro")
	if err != nil {
		return fmt.Errorf("failed to list ravro files: %w", err)
	}

	for _, file := range ravroFiles {
		if IsStructuredDataFile(file) {
			continue
		}

		// This is an attachment file - decrypt it
		decrypted, err := s.decryptFile(file, privateKeyPath)
		if err != nil {
			// Skip files that can't be decrypted
			continue
		}

		// Determine output path
		// Extract relative path from report path
		absReportPath, _ := filepath.Abs(reportPath)
		absFilePath, _ := filepath.Abs(file)

		relPath, err := filepath.Rel(absReportPath, absFilePath)
		if err != nil {
			// Fallback to basename
			relPath = filepath.Base(file)
		}

		// Remove .ravro extension
		relPath = strings.TrimSuffix(relPath, ".ravro")

		outputFile := filepath.Join(outputPath, relPath)

		// Ensure directory exists
		outputDir := filepath.Dir(outputFile)
		if err := s.storage.CreateDir(outputDir); err != nil {
			continue
		}

		// Write decrypted file
		if err := s.storage.WriteFile(outputFile, decrypted); err != nil {
			continue
		}
	}

	return nil
}

// CopyReportSidecars copies plaintext ZIP files needed for HTML (e.g. sha1_checksum.txt).
func (s *Service) CopyReportSidecars(reportPath string, outputPath string) error {
	for _, name := range []string{"sha1_checksum.txt"} {
		candidates := []string{
			filepath.Join(reportPath, name),
			filepath.Join(reportPath, "report", name),
		}
		for _, src := range candidates {
			if !s.storage.FileExists(src) {
				continue
			}
			if err := s.storage.CopyFile(src, filepath.Join(outputPath, name)); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// GetReportID extracts report ID from file path
func (s *Service) GetReportID(path string) string {
	// Extract report ID from path (e.g., encrypt/ir2020-07-16-0002)
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, part := range parts {
		if part == "encrypt" && i+1 < len(parts) {
			filename := parts[i+1]
			// Remove file extension if present
			filename = strings.TrimSuffix(filename, filepath.Ext(filename))
			return strings.TrimPrefix(filename, "report-")
		}
	}

	// Try to extract from filename
	base := filepath.Base(path)
	// Remove file extension if present
	base = strings.TrimSuffix(base, filepath.Ext(base))
	if strings.HasPrefix(base, "report-") {
		return strings.TrimPrefix(base, "report-")
	}

	return base
}
