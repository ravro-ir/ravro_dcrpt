package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ravro_dcrpt/internal/adapters/crypto"
	"ravro_dcrpt/internal/adapters/pdfgen"
	"ravro_dcrpt/internal/adapters/storage"
	"ravro_dcrpt/internal/core/decrypt"
	"ravro_dcrpt/internal/core/report"
)

const version = "v2.0.0"

type GUI struct {
	app         fyne.App
	window      fyne.Window
	inputEntry  *widget.Entry
	outputEntry *widget.Entry
	keyEntry    *widget.Entry
	logText     *widget.Entry
	progressBar *widget.ProgressBar
	processBtn  *widget.Button

	storageService *storage.FileSystemService
	cryptoService  crypto.PKCS7Service
	pdfGenerator   pdfgen.MarotoGenerator
	decryptService *decrypt.Service
	reportService  *report.Service
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Ravro Decryption Tool - " + version)

	gui := &GUI{
		app:    myApp,
		window: myWindow,
	}

	gui.setupServices()
	gui.setupUI()

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func (g *GUI) setupServices() {
	storageService := storage.NewFileSystemService()
	cryptoService := crypto.NewOpenSSLCGOService() // Using OpenSSL CGO for native PKCS7 decryption
	pdfGenerator := pdfgen.NewWKHTMLToPDFGenerator()

	decryptService := decrypt.NewService(cryptoService, storageService)
	reportService := report.NewService(decryptService, pdfGenerator, storageService)

	g.storageService = storageService.(*storage.FileSystemService)
	g.decryptService = decryptService
	g.reportService = reportService
}

func (g *GUI) setupUI() {
	// Title
	title := widget.NewLabelWithStyle(
		"🔐 Ravro Report Decryption Tool",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Input directory selection
	g.inputEntry = widget.NewEntry()
	g.inputEntry.SetPlaceHolder("Select input directory...")
	g.inputEntry.SetText("encrypt")

	inputBtn := widget.NewButton("Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				g.inputEntry.SetText(uri.Path())
			}
		}, g.window)
	})

	inputContainer := container.NewBorder(nil, nil, nil, inputBtn, g.inputEntry)

	// Output directory selection
	g.outputEntry = widget.NewEntry()
	g.outputEntry.SetPlaceHolder("Select output directory...")
	g.outputEntry.SetText("decrypt")

	outputBtn := widget.NewButton("Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				g.outputEntry.SetText(uri.Path())
			}
		}, g.window)
	})

	outputContainer := container.NewBorder(nil, nil, nil, outputBtn, g.outputEntry)

	// Key file selection
	g.keyEntry = widget.NewEntry()
	g.keyEntry.SetPlaceHolder("Select private key file...")
	g.keyEntry.SetText("key")

	keyBtn := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err == nil && uri != nil {
				g.keyEntry.SetText(uri.URI().Path())
			}
		}, g.window)
	})

	keyContainer := container.NewBorder(nil, nil, nil, keyBtn, g.keyEntry)

	// Form
	form := widget.NewForm(
		widget.NewFormItem("Input Directory", inputContainer),
		widget.NewFormItem("Output Directory", outputContainer),
		widget.NewFormItem("Private Key", keyContainer),
	)

	// Progress bar
	g.progressBar = widget.NewProgressBar()
	g.progressBar.Hide()

	// Log/Output area
	g.logText = widget.NewMultiLineEntry()
	g.logText.SetPlaceHolder("Output will appear here...")
	g.logText.Disable()

	logScroll := container.NewScroll(g.logText)
	logScroll.SetMinSize(fyne.NewSize(0, 200))

	// Buttons
	g.processBtn = widget.NewButton("🚀 Start Processing", g.processReports)

	initBtn := widget.NewButton("📁 Initialize Directories", g.initDirectories)

	validateBtn := widget.NewButton("🔍 Validate Key", g.validateKey)

	buttonContainer := container.NewHBox(
		g.processBtn,
		initBtn,
		validateBtn,
	)

	// Main content
	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			form,
			g.progressBar,
		),
		container.NewVBox(
			widget.NewSeparator(),
			buttonContainer,
		),
		nil,
		nil,
		logScroll,
	)

	g.window.SetContent(content)
}

func (g *GUI) log(message string) {
	current := g.logText.Text
	if current != "" {
		current += "\n"
	}
	g.logText.SetText(current + message)
	g.logText.CursorRow = len(g.logText.Text)
}

func (g *GUI) clearLog() {
	g.logText.SetText("")
}

func (g *GUI) initDirectories() {
	g.clearLog()
	g.log("📁 Initializing directories...")

	dirs := []string{"encrypt", "decrypt", "key"}
	for _, dir := range dirs {
		if err := g.storageService.CreateDir(dir); err != nil {
			g.log(fmt.Sprintf("❌ Failed to create %s: %v", dir, err))
			continue
		}
		g.log(fmt.Sprintf("✅ Created: %s/", dir))
	}

	g.log("\n✨ Directories initialized successfully!")
	g.log("📝 Next steps:")
	g.log("   1. Place your private key in the 'key/' directory")
	g.log("   2. Place encrypted reports (.zip or .ravro) in 'encrypt/' directory")
	g.log("   3. Click 'Start Processing' to decrypt and generate PDFs")

	dialog.ShowInformation("Success", "Directories initialized successfully!", g.window)
}

func (g *GUI) validateKey() {
	keyPath := g.keyEntry.Text
	if keyPath == "" {
		dialog.ShowError(fmt.Errorf("Please select a key file"), g.window)
		return
	}

	g.clearLog()
	g.log(fmt.Sprintf("🔍 Validating key: %s", keyPath))

	// Find actual key file
	actualKeyPath, err := g.findKey(keyPath)
	if err != nil {
		g.log(fmt.Sprintf("❌ Error: %v", err))
		dialog.ShowError(err, g.window)
		return
	}

	// Validate key
	if err := g.decryptService.ValidateKey(actualKeyPath); err != nil {
		g.log(fmt.Sprintf("❌ Invalid key: %v", err))
		dialog.ShowError(fmt.Errorf("Invalid key: %v", err), g.window)
		return
	}

	g.log("✅ Key is valid!")
	dialog.ShowInformation("Success", "Key is valid!", g.window)
}

func (g *GUI) processReports() {
	inputDir := g.inputEntry.Text
	outputDir := g.outputEntry.Text
	keyPath := g.keyEntry.Text

	if inputDir == "" || outputDir == "" || keyPath == "" {
		dialog.ShowError(fmt.Errorf("Please fill all fields"), g.window)
		return
	}

	// Disable button during processing
	g.processBtn.Disable()
	defer g.processBtn.Enable()

	g.clearLog()
	g.progressBar.Show()
	g.log("🚀 Starting processing...")

	// Find actual key file
	actualKeyPath, err := g.findKey(keyPath)
	if err != nil {
		g.log(fmt.Sprintf("❌ Error: %v", err))
		dialog.ShowError(err, g.window)
		g.progressBar.Hide()
		return
	}

	g.log(fmt.Sprintf("🔑 Using key: %s", actualKeyPath))

	// Validate key
	if err := g.decryptService.ValidateKey(actualKeyPath); err != nil {
		g.log(fmt.Sprintf("❌ Invalid key: %v", err))
		dialog.ShowError(fmt.Errorf("Invalid key: %v", err), g.window)
		g.progressBar.Hide()
		return
	}

	g.log(fmt.Sprintf("📂 Processing reports from: %s", inputDir))
	g.log(fmt.Sprintf("💾 Output directory: %s\n", outputDir))

	// Process reports
	results, err := g.reportService.ProcessReports(inputDir, actualKeyPath, outputDir)
	if err != nil {
		g.log(fmt.Sprintf("❌ Error: %v", err))
		dialog.ShowError(err, g.window)
		g.progressBar.Hide()
		return
	}

	// Display results
	g.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	g.log(fmt.Sprintf("%-30s %-15s %s", "Report ID", "Hunter", "Status"))
	g.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

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

		g.log(fmt.Sprintf("%-30s %-15s %s", result.ReportID, hunter, status))
	}

	g.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	g.log(fmt.Sprintf("\n📊 Summary: %d successful, %d failed, %d total",
		successCount, failCount, len(results)))

	if successCount > 0 {
		g.log(fmt.Sprintf("✨ PDFs generated in: %s", outputDir))
	}

	g.progressBar.Hide()

	// Show completion dialog
	if failCount == 0 {
		dialog.ShowInformation("Success",
			fmt.Sprintf("Successfully processed %d reports!", successCount),
			g.window)
	} else {
		dialog.ShowInformation("Completed with errors",
			fmt.Sprintf("Processed %d reports successfully, %d failed.", successCount, failCount),
			g.window)
	}
}

func (g *GUI) findKey(keyPath string) (string, error) {
	// Check if keyPath is a file
	if g.storageService.FileExists(keyPath) {
		return keyPath, nil
	}

	// Check if keyPath is a directory
	keys, err := g.storageService.ListFiles(keyPath, "*.pem")
	if err != nil || len(keys) == 0 {
		// Try any file in the directory
		keys, err = g.storageService.ListFiles(keyPath, "*")
		if err != nil {
			return "", fmt.Errorf("failed to find keys in %s: %w", keyPath, err)
		}

		// Filter out directories
		var files []string
		for _, key := range keys {
			if g.storageService.FileExists(key) {
				files = append(files, key)
			}
		}
		keys = files
	}

	if len(keys) == 0 {
		return "", fmt.Errorf("no key files found in %s", keyPath)
	}

	if len(keys) == 1 {
		return keys[0], nil
	}

	// Multiple keys found, use the first one
	g.log(fmt.Sprintf("⚠️  Multiple keys found, using: %s", filepath.Base(keys[0])))
	return keys[0], nil
}
