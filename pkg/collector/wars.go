package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchWars(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["wars"]; !exists {
		cc.cache["wars"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wars",
				Help: "Wars",
			},
			[]string{"war_id", "aggressor_name", "defender_name", "started", "finished", "mutual", "open_for_allies"},
		)
	}

	wars, _, err := cc.characters[characterID].ESIClient.Client.ESI.WarsApi.GetWars(cc.characters[characterID].ESIClient.Ctx, nil)
	if err != nil {
		return err
	}

	for _, war := range wars {
		localwar, _, err := cc.characters[characterID].ESIClient.Client.ESI.WarsApi.GetWarsWarId(cc.characters[characterID].ESIClient.Ctx, int32(war), nil)
		if err != nil {
			return err
		}
		cc.cache["wars"].WithLabelValues(strconv.Itoa(int(localwar.Id)), strconv.Itoa(int(localwar.Aggressor.AllianceId)), strconv.Itoa(int(localwar.Defender.AllianceId)), localwar.Declared.String(), localwar.Finished.String(), strconv.FormatBool(localwar.Mutual), strconv.FormatBool(localwar.OpenForAllies)).Set(1)
	}

	return nil
}
