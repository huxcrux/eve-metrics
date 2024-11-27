package helpers

import "github.com/huxcrux/eve-metrics/pkg/models"

func GetShipName(shipID int, esiClient models.ESIClient) string {
	// Lookup blueprint name
	ship, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseTypesTypeId(esiClient.Ctx, int32(shipID), nil)
	if err != nil {
		panic(err)
	}
	return ship.Name
}
