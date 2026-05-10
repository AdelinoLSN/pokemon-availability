package ports

type PokemonAvailabilityDetailExporter interface {
	ExportPokemonAvailabilityDetails(path string, data [][]string) error
}
