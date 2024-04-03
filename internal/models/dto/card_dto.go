package dto

// Represents one card
type CardDto struct {
	Suit  string `json:"suit"`  // The suit of this card
	Value string `json:"value"` // The value (full name) of this card
	Code  string `json:"code"`  // the code of this card
}
