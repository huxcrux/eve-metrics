package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	filePath = "config.yml"
)

type Config struct {
	Characters     []CharacterInput  `yaml:"characters"`
	Notification   NotificationInput `yaml:"notifications"`
	Discordwebhook string            `yaml:"discordwebhook"`
	ProxyURL       string            `yaml:"proxyurl"`
	Webhooks       []Webhook         `yaml:"webhooks"`
}

type Webhook struct {
	URL                      string  `yaml:"url"`
	AllAllainceSubscriptions bool    `yaml:"all_allaince_subscriptions"`
	AllianceSubscriptions    []int32 `yaml:"alliance_subscriptions"`
}
type CharacterInput struct {
	ID    int    `yaml:"id"`
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
}

type NotificationInput struct {
	Alliances    []NotificationAlliancesInput  `yaml:"alliances"`
	Corporations []NotificationCoporationInput `yaml:"corporations"`
	Characters   []NotificationCharacterInput  `yaml:"characters"`
}

type NotificationAlliancesInput struct {
	Character int32 `yaml:"character_id"`
	ID        int32 `yaml:"id"`
}

type NotificationCoporationInput struct {
	Character int32 `yaml:"character_id"`
	ID        int32 `yaml:"id"`
}

type NotificationCharacterInput struct {
	ID    int32  `yaml:"id"`
	Token string `yaml:"token"`
}

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
