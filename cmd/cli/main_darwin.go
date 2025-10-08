//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os"

	"ravro_dcrpt/internal/adapters/crypto"
	"ravro_dcrpt/internal/adapters/pdfgen"
	"ravro_dcrpt/internal/adapters/storage"
	"ravro_dcrpt/internal/core/decrypt"
	"ravro_dcrpt/internal/core/report"
	"ravro_dcrpt/internal/ports"

	"github.com/spf13/cobra"
)

const version = "v2.0.0-darwin"

var (
	inputDir  string
	outputDir string
	keyPath   string
	initDirs  bool
	jsonMode  bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ravro_dcrpt",
	Short: "Ravro Report Decryption Tool (macOS)",
	Long: `A versatile tool for decrypting and converting Ravro platform bug bounty reports to PDF.
macOS optimized version with OpenSSL support via Homebrew.`,
	Version: version,
	Run:     runDecrypt,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputDir, "in", "i", "encrypt", "Input directory containing encrypted reports")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "out", "o", "decrypt", "Output directory for decrypted reports")
	rootCmd.PersistentFlags().StringVarP(&keyPath, "key", "k", "key", "Path to private key file or directory")
	rootCmd.PersistentFlags().BoolVar(&initDirs, "init", false, "Initialize required directories")
	rootCmd.PersistentFlags().BoolVar(&jsonMode, "json", false, "Export reports as JSON")
}

func runDecrypt(cmd *cobra.Command, args []string) {
	// Initialize directories if requested
	if initDirs {
		initializeDirectories()
		return
	}

	// Display header
	displayHeader()

	// Validate directories
	if err := validateDirectories(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		fmt.Println("💡 Hint: Run with --init to create required directories")
		os.Exit(1)
	}

	// Initialize services
	storageService := storage.NewFileSystemService()

	// Use OpenSSL CGO service for macOS
	cryptoService := crypto.NewOpenSSLCGOService()
	pdfGenerator := pdfgen.NewWKHTMLToPDFGenerator()

	decryptService := decrypt.NewService(cryptoService, storageService)
	reportService := report.NewService(decryptService, pdfGenerator, storageService)

	// Find and validate key
	actualKeyPath, err := findKey(storageService, keyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔑 Using key: %s\n", actualKeyPath)

	// Validate key
	if err := decryptService.ValidateKey(actualKeyPath); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Invalid key: %v\n", err)
		os.Exit(1)
	}

	// Process reports
	fmt.Printf("📂 Processing reports from: %s\n", inputDir)
	fmt.Printf("💾 Output directory: %s\n\n", outputDir)

	results, err := reportService.ProcessReports(inputDir, actualKeyPath, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error processing reports: %v\n", err)
		os.Exit(1)
	}

	// Display results
	displayResults(results)
}

func initializeDirectories() {
	dirs := []string{"encrypt", "decrypt", "key"}
	storageService := storage.NewFileSystemService()

	for _, dir := range dirs {
		if err := storageService.CreateDir(dir); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to create %s: %v\n", dir, err)
			continue
		}
		fmt.Printf("✅ Created: %s/\n", dir)
	}

	fmt.Println("\n✨ Directories initialized successfully!")
	fmt.Println("📝 Next steps:")
	fmt.Println("   1. Place your private key in the 'key/' directory")
	fmt.Println("   2. Place encrypted reports in the 'encrypt/' directory")
	fmt.Println("   3. Run: ravro_dcrpt")
	fmt.Println("\n🍺 macOS Note: Make sure OpenSSL is installed via Homebrew:")
	fmt.Println("   brew install openssl")
}

func validateDirectories() error {
	storageService := storage.NewFileSystemService()

	requiredDirs := []string{"encrypt", "decrypt", "key"}
	for _, dir := range requiredDirs {
		if !storageService.FileExists(dir) {
			return fmt.Errorf("required directory not found: %s", dir)
		}
	}

	return nil
}

func findKey(storageService ports.StorageService, keyPath string) (string, error) {
	// Check if keyPath is a file
	if storageService.FileExists(keyPath) {
		return keyPath, nil
	}

	// Check if keyPath is a directory
	keys, err := storageService.ListFiles(keyPath, "*.pem")
	if err != nil {
		keys, err = storageService.ListFiles(keyPath, "*")
		if err != nil {
			return "", fmt.Errorf("failed to find keys in %s: %w", keyPath, err)
		}
	}

	if len(keys) == 0 {
		return "", fmt.Errorf("no key files found in %s", keyPath)
	}

	if len(keys) == 1 {
		return keys[0], nil
	}

	// Multiple keys found, use the first one (or implement selection)
	fmt.Printf("⚠️  Multiple keys found, using: %s\n", keys[0])
	return keys[0], nil
}

func displayHeader() {
	fmt.Println("╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║    Ravro Report Decryption Tool - macOS Version      ║")
	fmt.Printf("║                  Version: %-24s  ║\n", version)
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	fmt.Println()
}

func displayResults(results []*report.ProcessResult) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("%-30s %-20s %s\n", "Report ID", "Hunter", "Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	successCount := 0
	failCount := 0

	for _, result := range results {
		status := "✅ Success"
		if !result.Success {
			status = "❌ Failed"
			failCount++
		} else {
			successCount++
		}

		hunter := "N/A"
		if result.Report != nil {
			hunter = result.Report.HunterUsername
		}

		fmt.Printf("%-30s %-20s %s\n", result.ReportID, hunter, status)

		// Show error details for failed reports
		if !result.Success && result.Error != nil {
			fmt.Printf("   Error: %v\n", result.Error)
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n📊 Summary: %d successful, %d failed, %d total\n",
		successCount, failCount, len(results))

	if successCount > 0 {
		fmt.Printf("✨ PDFs generated in: %s\n", outputDir)
	}
}
