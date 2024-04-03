package models

import "github.com/google/uuid"

// A deck of cards
type Deck struct {
	DeckId    uuid.UUID // The deck's Id
	Cards     [][2]rune // Card code
	Shuffled  bool      // Indicates if the deck has been shuffled
	Remaining uint8     // Number of remaining cards
}
