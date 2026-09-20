package csvreport

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

func TestAvailabilityWriterWritesReportAsCSV(t *testing.T) {
	outputDirectory := t.TempDir()
	writer := NewAvailabilityWriter(outputDirectory)
	err := writer.Write(context.Background(), exportavailability.Report{
		Sequence: 2,
		Game:     domain.Game{Abbreviation: "R"},
		Rows: []exportavailability.ReportRow{{
			Number: "0001", Name: "Bulbasaur", Methods: "WILD",
		}},
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDirectory, "002_R.csv"))
	if err != nil {
		t.Fatal(err)
	}
	want := "Number,Name,Form,Methods,Notes\n0001,Bulbasaur,,WILD,\n"
	if string(data) != want {
		t.Fatalf("CSV = %q, want %q", data, want)
	}
}
