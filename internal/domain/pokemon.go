package domain

type Pokemon struct {
	Number       int                   `json:"number"`
	Name         string                `json:"name"`
	Form         string                `json:"form"`
	Availability []PokemonAvailability `json:"availability"`
}

type PokemonAvailability struct {
	Game   string  `json:"game"`
	Method string  `json:"method"`
	Notes  *string `json:"notes,omitempty"`
}
