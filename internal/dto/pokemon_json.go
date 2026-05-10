package dto

type PokemonJson struct {
	Number       int                       `json:"number"`
	Name         string                    `json:"name"`
	Form         string                    `json:"form"`
	Availability []PokemonAvailabilityJson `json:"availability"`
}
