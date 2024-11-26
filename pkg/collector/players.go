package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchPlayersOnline(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["players_online"]; !exists {
		cc.cache["players_online"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "players_online",
				Help: "Number of players online",
			},
			[]string{},
		)
	}

	// Fetch player count
	players, _, err := cc.characters[characterID].ESIClient.Client.ESI.StatusApi.GetStatus(cc.characters[characterID].ESIClient.Ctx, nil)
	if err != nil {
		return err
	}

	cc.cache["players_online"].WithLabelValues().Set(float64(players.Players))

	return nil
}
