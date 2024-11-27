package helpers

import (
	"fmt"

	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
	"github.com/huxcrux/eve-metrics/pkg/models"
)

func GetCharacterName(id int32, esiClient models.ESIClient) string {

	// Get character name from id
	character, _, err := esiClient.Client.ESI.CharacterApi.GetCharactersCharacterId(esiClient.Ctx, id, nil)
	if err != nil {
		fmt.Println("Error fetching character")
		return "Unknown"
	}

	return character.Name

}

func GetCharacterInfo(id int32, token string) (models.Character, error) {

	esiClient := esiClient.NewESIClient(token)
	// Get character info from id
	character, _, err := esiClient.Client.ESI.CharacterApi.GetCharactersCharacterId(esiClient.Ctx, id, nil)
	if err != nil {
		fmt.Println("Error fetching character")
		return models.Character{}, err
	}

	characterInfo := models.Character{
		ID:            int(id),
		CorporationID: int(character.CorporationId),
		AllianceID:    int(character.AllianceId),
		Name:          character.Name,
		ESIClient:     esiClient,
	}
	return characterInfo, nil

}
