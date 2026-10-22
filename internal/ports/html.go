package ports

import "ravro_dcrpt/pkg/models"

// HTMLGenerator defines the interface for standalone HTML report export.
type HTMLGenerator interface {
	GenerateReport(report *models.Report, judgment *models.Judgment, amendments []models.Amendment, outputPath string) error
}
