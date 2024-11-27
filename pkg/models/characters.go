package models

type Character struct {
	ID            int
	CorporationID int
	AllianceID    int
	Name          string
	ESIClient     ESIClient
}
