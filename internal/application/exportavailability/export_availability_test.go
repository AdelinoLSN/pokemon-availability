package exportavailability_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

func TestUseCaseBuildsSemanticReport(t *testing.T) {
	reader := &readerStub{
		games: []domain.Game{{Abbreviation: "R"}},
		rows: map[string][]exportavailability.AvailabilityRow{
			"R": {
				{Number: 25, Name: "Pikachu", MethodKey: "STARTER", Note: "Only, one"},
				{Number: 25, Name: "Pikachu", MethodKey: "WILD"},
			},
		},
	}
	writer := &writerSpy{}

	err := exportavailability.New(reader, writer).Execute(context.Background(), exportavailability.Command{})

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	row := writer.reports[0].Rows[0]
	if row.Number != "0025" || row.Methods != "STARTER|WILD" || row.Notes != "[STARTER]: Only; one|" {
		t.Fatalf("row = %#v", row)
	}
}

func TestUseCaseDoesNotWriteWhenReaderFails(t *testing.T) {
	expected := errors.New("read failed")
	writer := &writerSpy{}
	reader := &readerStub{games: []domain.Game{{Abbreviation: "R"}}, err: expected}

	err := exportavailability.New(reader, writer).Execute(context.Background(), exportavailability.Command{})

	if !errors.Is(err, expected) {
		t.Fatalf("Execute() error = %v, want %v", err, expected)
	}
	if len(writer.reports) != 0 {
		t.Fatal("writer must not be called")
	}
}

type readerStub struct {
	games []domain.Game
	rows  map[string][]exportavailability.AvailabilityRow
	err   error
}

func (s *readerStub) ListGames(context.Context) ([]domain.Game, error) { return s.games, nil }
func (s *readerStub) ByGame(_ context.Context, game string) ([]exportavailability.AvailabilityRow, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.rows[game], nil
}

type writerSpy struct {
	reports []exportavailability.Report
}

func (s *writerSpy) Write(_ context.Context, report exportavailability.Report) error {
	s.reports = append(s.reports, report)
	return nil
}
