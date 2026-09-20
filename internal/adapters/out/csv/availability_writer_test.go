package csvreport

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

func TestAvailabilityWriterRendersReportAsCSV(t *testing.T) {
	output := t.TempDir()
	err := NewAvailabilityWriter(output).Write(context.Background(), exportavailability.Report{
		Sequence: 2,
		Game:     domain.Game{Abbreviation: "R"},
		Rows:     []exportavailability.ReportRow{{Number: "0001", Name: "Bulbasaur", Methods: "WILD"}},
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(output, "002_R.csv"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(data), "Number,Name,Form,Methods,Notes\n0001,Bulbasaur,,WILD,\n"; got != want {
		t.Fatalf("CSV = %q, want %q", got, want)
	}
}
