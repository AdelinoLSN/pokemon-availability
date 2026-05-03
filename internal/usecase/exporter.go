package usecase

import (
  "database/sql"
  "fmt"
  "log"
  "strings"

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

    rows, err := database.GetPokemonFullViewForGame(e.db, game.Abbreviation)
    if err != nil {
      return fmt.Errorf("falha ao buscar dados do jogo %s: %w", game.Abbreviation, err)
    }

    csvData := e.mapToCSVFormat(rows)

    if err := e.csvHandler.WriteAll(filename, csvData); err != nil {
      return fmt.Errorf("falha ao escrever CSV do jogo %s: %w", game.Abbreviation, err)
    }

    // log.Printf("Exportado com sucesso: %s (%d registros agrupados)", filename, len(csvData)-1)
  }

  return nil
}

func (e *Exporter) mapToCSVFormat(rows []database.PokemonFullViewRow) [][]string {
  var data [][]string
  data = append(data, []string{"Number", "Name", "Form", "Methods", "Notes"})

  if len(rows) == 0 {
    return data
  }

  var currentNumber int
  var currentName, currentForm string
  var currentMethods []string
  var currentNotes []string

  appendCurrentPokemon := func() {
    methodsStr := strings.Join(currentMethods, "|")

    notesStr := strings.Join(currentNotes, "|")
    notesStr = strings.ReplaceAll(notesStr, ",", ";")

    data = append(data, []string{
      fmt.Sprintf("%04d", currentNumber),
      currentName,
      currentForm,
      methodsStr,
      notesStr,
    })
  }

  for i, row := range rows {
    if i == 0 || row.Number != currentNumber || row.Form != currentForm {
      if i > 0 {
        appendCurrentPokemon()
      }

      currentNumber = row.Number
      currentName = row.Name
      currentForm = row.Form
      currentMethods = []string{}
      currentNotes = []string{}
    }

    currentMethods = append(currentMethods, row.Method)
    currentNotes = append(currentNotes, row.Note)
  }

  if len(rows) > 0 {
    appendCurrentPokemon()
  }

  return data
}
