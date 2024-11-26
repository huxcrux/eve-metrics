package collector

import (
	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchJumpFatigue(characterID int) error {
	// Fetch all corporations in an alliance

	if _, exists := cc.cache["character_jump_fatigue"]; !exists {
		cc.cache["character_jump_fatigue"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character_jump_fatigue",
				Help: "Character Jump Fatigue",
			},
			[]string{"character_name", "last_jump_date", "last_update_date"},
		)
	}

	jumpFatigue, _, err := cc.characters[characterID].ESIClient.Client.ESI.CharacterApi.GetCharactersCharacterIdFatigue(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	cc.cache["character_jump_fatigue"].WithLabelValues(characterName, jumpFatigue.LastJumpDate.String(), jumpFatigue.LastUpdateDate.String()).Set(float64(jumpFatigue.JumpFatigueExpireDate.Unix()))

	return nil
}
