package domain

type Game struct {
	Abbreviation string `json:"abbreviation"`
	Name         string `json:"name"`
	Generation   int    `json:"generation"`
}
