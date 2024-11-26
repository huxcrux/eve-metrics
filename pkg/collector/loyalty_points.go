package collector

import (
	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchCharacterLoyaltyPoints(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["character_loyalty_points"]; !exists {
		cc.cache["character_loyalty_points"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character_loyalty_points",
				Help: "Character loyalty points",
			},
			[]string{"character_name", "corporation_name"},
		)
	}

	// Fetch player count
	loyaltyPoints, _, err := cc.characters[characterID].ESIClient.Client.ESI.LoyaltyApi.GetCharactersCharacterIdLoyaltyPoints(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	for _, loyaltyPoint := range loyaltyPoints {
		coporationName := helpers.GetCorporationName(loyaltyPoint.CorporationId, cc.characters[characterID].ESIClient)
		cc.cache["character_loyalty_points"].WithLabelValues(characterName, coporationName).Set(float64(loyaltyPoint.LoyaltyPoints))
	}

	return nil
}
