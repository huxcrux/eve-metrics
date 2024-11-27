package collector

import (
	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchAllianceCorporations(characterID int) error {
	// Fetch all corporations in an alliance

	if _, exists := cc.cache["alliance_corporation"]; !exists {
		cc.cache["alliance_corporation"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "alliance_corporation",
				Help: "Corporations in an alliance",
			},
			[]string{"alliance_name", "corporation_name"},
		)
	}

	corporations, _, err := cc.characters[characterID].ESIClient.Client.ESI.AllianceApi.GetAlliancesAllianceIdCorporations(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].AllianceID), nil)
	if err != nil {
		return err
	}

	allianceName := helpers.GetAllianceName(int32(cc.characters[characterID].AllianceID), cc.characters[characterID].ESIClient)

	for _, corporation := range corporations {
		corpName := helpers.GetCorporationName(corporation, cc.characters[characterID].ESIClient)
		cc.cache["alliance_corporation"].WithLabelValues(allianceName, corpName).Set(1)
	}

	return nil
}
