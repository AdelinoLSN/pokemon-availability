package exportavailability

import (
	"context"
	"fmt"
	"strings"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

type Command struct{}

type InputPort interface {
	Execute(context.Context, Command) error
}

type Reader interface {
	ListGames(context.Context) ([]domain.Game, error)
	ByGame(context.Context, string) ([]AvailabilityRow, error)
}

type Writer interface {
	Write(context.Context, Report) error
}

type AvailabilityRow struct {
	Number    int
	Name      string
	Form      string
	MethodKey string
	Note      string
}

type Report struct {
	Sequence int
	Game     domain.Game
	Rows     []ReportRow
}

type ReportRow struct {
	Number  string
	Name    string
	Form    string
	Methods string
	Notes   string
}

type UseCase struct {
	reader Reader
	writer Writer
}

func New(reader Reader, writer Writer) *UseCase {
	return &UseCase{reader: reader, writer: writer}
}

func (u *UseCase) Execute(ctx context.Context, _ Command) error {
	games, err := u.reader.ListGames(ctx)
	if err != nil {
		return err
	}

	for sequence, game := range games {
		rows, err := u.reader.ByGame(ctx, game.Abbreviation)
		if err != nil {
			return err
		}

		report := Report{
			Sequence: sequence,
			Game:     game,
			Rows:     buildReportRows(rows),
		}
		if err := u.writer.Write(ctx, report); err != nil {
			return err
		}
	}

	return nil
}

func buildReportRows(rows []AvailabilityRow) []ReportRow {
	if len(rows) == 0 {
		return []ReportRow{}
	}

	reportRows := make([]ReportRow, 0)
	current := rows[0]
	methods := make([]string, 0)
	notes := make([]string, 0)

	appendCurrent := func() {
		reportRows = append(reportRows, ReportRow{
			Number:  fmt.Sprintf("%04d", current.Number),
			Name:    current.Name,
			Form:    current.Form,
			Methods: strings.Join(methods, "|"),
			Notes:   strings.ReplaceAll(strings.Join(notes, "|"), ",", ";"),
		})
	}

	for index, row := range rows {
		isNewPokemon := index > 0 &&
			(row.Number != current.Number || row.Form != current.Form)
		if isNewPokemon {
			appendCurrent()
			current = row
			methods = methods[:0]
			notes = notes[:0]
		}

		methods = append(methods, row.MethodKey)
		note := row.Note
		if note != "" {
			note = fmt.Sprintf("[%s]: %s", row.MethodKey, note)
		}
		notes = append(notes, note)
	}

	appendCurrent()
	return reportRows
}
