package configs

// All configuration data for api and decks service
type Config struct {
	Api   ApiConfig   `yaml:"api"`
	Decks DecksConfig `yaml:"decks"`
}
