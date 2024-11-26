package helpers

import (
	"fmt"

	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
)

func GetStructureName(structureID int64, esiClient esiClient.ESIClient) string {
	// Lookup structure name
	structure, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseStructuresStructureId(esiClient.Ctx, int64(structureID), nil)
	if err != nil {
		if err.Error() != "400 Bad Request" {
			panic(err)
		}
	} else {
		return structure.Name
	}

	station, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseStationsStationId(esiClient.Ctx, int32(structureID), nil)
	if err != nil {
		fmt.Println("Error fetching station details:", err)
	} else {
		return station.Name
	}

	return "Unknown"

}

func GetStructureSystem(structureID int64, esiClient esiClient.ESIClient) string {
	// Lookup structure system
	structure, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseStructuresStructureId(esiClient.Ctx, int64(structureID), nil)
	if err != nil {
		if err.Error() != "400 Bad Request" {
			panic(err)
		}
	} else {
		return GetSolarSystemName(structure.SolarSystemId, esiClient)
	}

	// check if structure is empty
	station, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseStationsStationId(esiClient.Ctx, int32(structureID), nil)
	if err != nil {
		fmt.Println("Error fetching station details:", err)
	} else {
		return GetSolarSystemName(station.SystemId, esiClient)
	}

	return "Unknown"

}

func GetStructureOwner(structureID int64, esiClient esiClient.ESIClient) string {
	// Lookup structure owner
	structure, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseStructuresStructureId(esiClient.Ctx, int64(structureID), nil)
	if err != nil {
		if err.Error() != "400 Bad Request" {
			panic(err)
		}
	} else {
		return GetCorporationName(structure.OwnerId, esiClient)
	}

	// check if structure is empty
	station, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseStationsStationId(esiClient.Ctx, int32(structureID), nil)
	if err != nil {
		fmt.Println("Error fetching station details:", err)
	} else {
		return GetCorporationName(station.Owner, esiClient)
	}

	return "Unknown"
}
