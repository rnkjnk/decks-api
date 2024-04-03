package configs

// Confiuration for the decks service
type DecksConfig struct {
	// The first character of the suits and values names will be taken to form codes of cards
	// That means that all suits and values respectively need to start with different letters
	Suits  []string `yaml:"suits"`  // list of names of suits
	Values []string `yaml:"values"` // list of names of values
}
