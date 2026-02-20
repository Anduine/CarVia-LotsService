package domain

type LotsResponse struct {
	Lots  []Lot `json:"lots"`
	Total int   `json:"total"`
}
