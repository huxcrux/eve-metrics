package helpers

import (
	"fmt"

	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
)

func GetAllianceName(id int32, esiClient esiClient.ESIClient) string {
	// Fetch owner alliance details
	alliance, _, err := esiClient.Client.ESI.AllianceApi.GetAlliancesAllianceId(esiClient.Ctx, id, nil)
	if err != nil {
		fmt.Println("Error fetching alliance details:", err)
		return ""
	}

	return alliance.Name
}
