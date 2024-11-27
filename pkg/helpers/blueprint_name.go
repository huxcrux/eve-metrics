package helpers

import "github.com/huxcrux/eve-metrics/pkg/models"

func GetBlueprintName(blueprintID int, esiClient models.ESIClient) string {
	// Lookup blueprint name
	blueprint, _, err := esiClient.Client.ESI.UniverseApi.GetUniverseTypesTypeId(esiClient.Ctx, int32(blueprintID), nil)
	if err != nil {
		panic(err)
	}
	return blueprint.Name
}
