package helpers

import esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"

func GetShipName(shipID int, esiClient esiClient.ESIClient) string {
	// Lookup blueprint name
	ship, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseTypesTypeId(esiClient.Ctx, int32(shipID), nil)
	if err != nil {
		panic(err)
	}
	return ship.Name
}
