package config

import (
	"fmt"
	"log"
	"os"

	"github.com/huxcrux/eve-metrics/pkg/models"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Characters []models.CharacterInput `yaml:"characters"`
}

func ReadConfig() []models.CharacterInput {
	// Path to your YAML file

	// Path to your YAML file
	filePath := "config.yml"

	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the YAML content into the Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Print the parsed characters
	for i, character := range config.Characters {
		fmt.Printf("Character %d: %+v\n", i+1, character)
	}

	return config.Characters
}
