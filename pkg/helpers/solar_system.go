package helpers

import (
	"fmt"

	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
)

func GetSolarSystemName(id int32, esiClient esiClient.ESIClient) string {

	// Fetch solar system details
	solarSystem, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseSystemsSystemId(esiClient.Ctx, id, nil)
	if err != nil {
		fmt.Println("Error fetching solar system details:", err)
		return ""
	}

	return solarSystem.Name

}
