package config

import (
	"log"
	"os"

	"github.com/huxcrux/eve-metrics/pkg/models"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Characters     []models.CharacterInput  `yaml:"characters"`
	Notification   models.NotificationInput `yaml:"notifications"`
	Discordwebhook string                   `yaml:"discordwebhook"`
}

var (
	filePath = "config.yml"
)

func ReadConfig() Config {

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

	return config
}
