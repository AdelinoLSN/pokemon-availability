package postgres

var Tables = struct {
	Games               string
	Methods             string
	Pokemon             string
	PokemonAvailability string
}{
	Games:               "games",
	Methods:             "methods",
	Pokemon:             "pokemon",
	PokemonAvailability: "pokemon_availability",
}

var Views = struct {
	PokemonAvailabilityDetails                string
	PokemonAvailabilityDetailsInconsistencies string
}{
	PokemonAvailabilityDetails:                "v_pokemon_availability_details",
	PokemonAvailabilityDetailsInconsistencies: "v_pokemon_availability_details_inconsistencies",
}

var MaterializedViews = struct {
	PokemonAvailabilityDetails string
}{
	PokemonAvailabilityDetails: "mv_pokemon_availability_details",
}
