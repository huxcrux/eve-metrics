package collector

import (
	"fmt"

	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchCharacterOnlineStatus(characterID int) error {
	// Fetch character online status

	if _, exists := cc.cache["character_online_status"]; !exists {
		cc.cache["character_online_status"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character_online_status",
				Help: "Character online status",
			},
			[]string{"character_name", "last_login", "last_logout", "logins"},
		)
	}

	onlineStatus, _, err := cc.characters[characterID].ESIClient.Client.ESI.LocationApi.GetCharactersCharacterIdOnline(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	if onlineStatus.Online {
		cc.cache["character_online_status"].WithLabelValues(characterName, onlineStatus.LastLogin.String(), onlineStatus.LastLogout.String(), fmt.Sprintf("%d", onlineStatus.Logins)).Set(1)
	} else {
		cc.cache["character_online_status"].WithLabelValues(characterName, onlineStatus.LastLogin.String(), onlineStatus.LastLogout.String(), fmt.Sprintf("%d", onlineStatus.Logins)).Set(0)
	}
	return nil
}

func (cc *CachedCollector) FetchCharacterLocation(characterID int) error {
	// Fetch character location

	if _, exists := cc.cache["character_location"]; !exists {
		cc.cache["character_location"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character_location",
				Help: "Character location",
			},
			[]string{"character_name", "solar_system_name", "structure_name"},
		)
	}

	location, _, err := cc.characters[characterID].ESIClient.Client.ESI.LocationApi.GetCharactersCharacterIdLocation(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	locationName := helpers.GetSolarSystemName(location.SolarSystemId, cc.characters[characterID].ESIClient)
	structureName := helpers.GetStructureName(location.StructureId, cc.characters[characterID].ESIClient)
	cc.cache["character_location"].WithLabelValues(characterName, locationName, structureName).Set(1)
	return nil
}

func (cc *CachedCollector) FetchCharacterShip(characterID int) error {
	// Fetch character ship

	if _, exists := cc.cache["character_ship"]; !exists {
		cc.cache["character_ship"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "character_ship",
				Help: "Character ship",
			},
			[]string{"character_name", "ship_name"},
		)
	}

	ship, _, err := cc.characters[characterID].ESIClient.Client.ESI.LocationApi.GetCharactersCharacterIdShip(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	shipName := helpers.GetShipName(int(ship.ShipTypeId), cc.characters[characterID].ESIClient)
	cc.cache["character_ship"].WithLabelValues(characterName, shipName).Set(1)
	return nil
}
