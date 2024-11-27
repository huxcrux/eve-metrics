package helpers

import (
	"fmt"

	"github.com/huxcrux/eve-metrics/pkg/models"
)

func GetCorporationName(id int32, esiClient models.ESIClient) string {
	// Fetch owner corporation details
	corporation, _, err := esiClient.Client.ESI.CorporationApi.GetCorporationsCorporationId(esiClient.Ctx, id, nil)
	if err != nil {
		fmt.Println("Error fetching corporation details:", err)
		return ""
	}

	return corporation.Name
}
