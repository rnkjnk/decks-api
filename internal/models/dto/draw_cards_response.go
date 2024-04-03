package dto

// DTO for deck object
type DrawCardsResponse struct {
	Cards []CardDto `json:"cards"` // The cards
}
