package helpers

import esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"

func GetBlueprintName(blueprintID int, esiClient esiClient.ESIClient) string {
	// Lookup blueprint name
	blueprint, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseTypesTypeId(esiClient.Ctx, int32(blueprintID), nil)
	if err != nil {
		panic(err)
	}
	return blueprint.Name
}
