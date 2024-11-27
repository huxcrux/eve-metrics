package models

import esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"

type Character struct {
	ID            int
	CorporationID int
	AllianceID    int
	Name          string
	ESIClient     esiClient.ESIClient
}
