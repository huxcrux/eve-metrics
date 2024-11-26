package helpers

import (
	"fmt"

	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
)

func GetCorporationName(id int32, esiClient esiClient.ESIClient) string {
	// Fetch owner corporation details
	corporation, _, err := esiClient.Client.ESI.CorporationApi.GetCorporationsCorporationId(esiClient.Ctx, id, nil)
	if err != nil {
		fmt.Println("Error fetching corporation details:", err)
		return ""
	}

	return corporation.Name
}
