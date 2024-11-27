package helpers

import (
	"fmt"

	"github.com/huxcrux/eve-metrics/pkg/models"
)

func GetAllianceName(id int32, esiClient models.ESIClient) string {
	// Fetch owner alliance details
	alliance, _, err := esiClient.Client.ESI.AllianceApi.GetAlliancesAllianceId(esiClient.Ctx, id, nil)
	if err != nil {
		fmt.Println("Error fetching alliance details:", err)
		return ""
	}

	return alliance.Name
}
