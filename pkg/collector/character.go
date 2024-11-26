package collector

import (
	"fmt"

	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchCharacter(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["character"]; !exists {
		cc.cache["character"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character",
				Help: "Character information",
			},
			[]string{"character_name", "corporation_name", "alliance_name", "security_status", "created_at"},
		)
	}

	// Fetch character
	character, _, err := cc.characters[characterID].ESIClient.Client.ESI.CharacterApi.GetCharactersCharacterId(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	corporationName := helpers.GetCorporationName(character.CorporationId, cc.characters[characterID].ESIClient)
	allianceName := helpers.GetAllianceName(character.AllianceId, cc.characters[characterID].ESIClient)

	cc.cache["character"].WithLabelValues(character.Name, corporationName, allianceName, fmt.Sprintf("%f", character.SecurityStatus), character.Birthday.String()).Set(1)

	return nil
}
