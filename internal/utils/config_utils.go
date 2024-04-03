package utils

import (
	"os"

	"github.com/rnkjnk/decks-api/internal/models/configs"
	"gopkg.in/yaml.v3"
)

func GetConfigsFromYaml(path string) configs.Config {
	// Read YAML file
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to read configuration from YAML file")
	}

	// Parse YAML into Config struct
	var config configs.Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic("Failed to read configuration from YAML file")
	}

	return config
}
