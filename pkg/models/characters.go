package models

import esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"

type CharacterInput struct {
	ID    int    `yaml:"id"`
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
}

type Character struct {
	ID            int
	CorporationID int
	AllianceID    int
	Name          string
	ESIClient     esiClient.ESIClient
}
