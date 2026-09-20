package csvreport

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/filesystem"
)

var _ exportavailability.Writer = (*AvailabilityWriter)(nil)

type AvailabilityWriter struct {
	outputDirectory string
}

func NewAvailabilityWriter(outputDirectory string) *AvailabilityWriter {
	return &AvailabilityWriter{outputDirectory: outputDirectory}
}

func (w *AvailabilityWriter) Write(ctx context.Context, report exportavailability.Report) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data := make([][]string, 0, len(report.Rows)+1)
	data = append(data, []string{"Number", "Name", "Form", "Methods", "Notes"})
	for _, row := range report.Rows {
		data = append(data, []string{row.Number, row.Name, row.Form, row.Methods, row.Notes})
	}
	filename := fmt.Sprintf("%03d_%s.csv", report.Sequence, report.Game.Abbreviation)
	return filesystem.WriteCSV(filepath.Join(w.outputDirectory, filename), data)
}
