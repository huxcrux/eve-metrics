package collector

import (
	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchCharacterWallet(characterID int) error {
	// Fetch all corporations in an alliance

	if _, exists := cc.cache["character_wallet"]; !exists {
		cc.cache["character_wallet"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character_wallet",
				Help: "Character wallet",
			},
			[]string{"character_name"},
		)
	}

	wallet, _, err := cc.characters[characterID].ESIClient.Client.ESI.WalletApi.GetCharactersCharacterIdWallet(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	cc.cache["character_wallet"].WithLabelValues(characterName).Set(wallet)

	return nil
}
