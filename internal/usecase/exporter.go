package usecase

import (
	"database/sql"
	"fmt"
	"log"

	"pokemon-availability/internal/domain"
	"pokemon-availability/internal/adapters/database"
)

type CSVHandler interface {
	WriteAll(path string, data [][]string) error
}

type Exporter struct {
	db         *sql.DB
	csvHandler CSVHandler
}

func NewExporter(db *sql.DB, csvHandler CSVHandler) *Exporter {
	return &Exporter{
		db:         db,
		csvHandler: csvHandler,
	}
}

func (e *Exporter) ExportGamesToCSV(games []domain.Game) error {
	log.Println("Iniciando exportação para CSV...")

	for i, game := range games {
		id := i + 1
		filename := fmt.Sprintf(".outputs/csv/%03d-%s.csv", id, game.Abbreviation)

		rows, err := database.GetPokemonForGame(e.db, game.Abbreviation)
		if err != nil {
			return fmt.Errorf("falha ao buscar dados do jogo %s: %w", game.Abbreviation, err)
		}

		csvData := e.mapToCSVFormat(rows)

		if err := e.csvHandler.WriteAll(filename, csvData); err != nil {
			return fmt.Errorf("falha ao escrever CSV do jogo %s: %w", game.Abbreviation, err)
		}

		log.Printf("Exportado com sucesso: %s (%d registros)", filename, len(rows))
	}

	return nil
}

func (e *Exporter) mapToCSVFormat(rows []database.PokemonCSVRow) [][]string {
	var data [][]string

	data = append(data, []string{"Number", "Name", "Form", "Methods", "Notes"})

	for _, row := range rows {
		data = append(data, []string{
			fmt.Sprintf("%04d", row.Number),
			row.Name,
			row.Form,
			row.Methods,
			row.Notes,
		})
	}

	return data
}
