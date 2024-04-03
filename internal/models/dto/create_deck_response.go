package dto

// DTO for deck object
type CreateDeckResponse struct {
	DeckId    string `json:"deck_id"`   // The Id of the deck (a uuid represented as string)
	Shuffled  bool   `json:"shuffled"`  // If the deck has been shuffled
	Remaining uint8  `json:"remaining"` // Number of remaining cards
}
