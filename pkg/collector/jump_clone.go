package collector

import (
	"fmt"

	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchCharacterJumpClones(characterID int) error {
	// Fetch all corporations in an alliance

	if _, exists := cc.cache["character_jump_clone"]; !exists {
		cc.cache["character_jump_clone"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character_jump_clone",
				Help: "Character Jump Clone",
			},
			[]string{"character_name", "home_location", "structure_name", "structure_owner", "name", "Implants"},
		)
	}

	jumpClone, _, err := cc.characters[characterID].ESIClient.Client.ESI.ClonesApi.GetCharactersCharacterIdClones(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	homeLocationName := helpers.GetStructureName(jumpClone.HomeLocation.LocationId, cc.characters[characterID].ESIClient)
	for _, clone := range jumpClone.JumpClones {
		locationName := helpers.GetStructureName(clone.LocationId, cc.characters[characterID].ESIClient)
		structureOwner := helpers.GetStructureOwner(clone.LocationId, cc.characters[characterID].ESIClient)
		cc.cache["character_jump_clone"].WithLabelValues(characterName, homeLocationName, locationName, structureOwner, clone.Name, fmt.Sprintf("%v", clone.Implants)).Set(1)
	}

	return nil
}
