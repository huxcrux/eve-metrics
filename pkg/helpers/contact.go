package helpers

import "github.com/huxcrux/eve-metrics/pkg/models"

func GetCharacterContactLabelName(lableID int64, charaterID int32, esiClient models.ESIClient) string {
	// Get contact label name
	contactLabel, _, err := esiClient.Client.ESI.ContactsApi.GetCharactersCharacterIdContactsLabels(esiClient.Ctx, charaterID, nil)
	if err != nil {
		return "Unknown"
	}

	for _, label := range contactLabel {
		if label.LabelId == lableID {
			return label.LabelName
		}
	}

	return "Unknown"
}

func GetAllianceContactLabelName(lableID int64, allianceID int32, esiClient models.ESIClient) string {
	// Get contact label name
	contactLabel, _, err := esiClient.Client.ESI.ContactsApi.GetAlliancesAllianceIdContactsLabels(esiClient.Ctx, allianceID, nil)
	if err != nil {
		return "Unknown"
	}

	for _, label := range contactLabel {
		if label.LabelId == lableID {
			return label.LabelName
		}
	}

	return "Unknown"
}

func GetCorporationContactLabelName(lableID int64, corporationID int32, esiClient models.ESIClient) string {
	// Get contact label name
	contactLabel, _, err := esiClient.Client.ESI.ContactsApi.GetCorporationsCorporationIdContactsLabels(esiClient.Ctx, corporationID, nil)
	if err != nil {
		return "Unknown"
	}

	for _, label := range contactLabel {
		if label.LabelId == lableID {
			return label.LabelName
		}
	}

	return "Unknown"
}

func GetContactName(contactType string, contactID int32, esiclient models.ESIClient) string {
	switch contactType {
	case "character":
		return GetCharacterName(contactID, esiclient)
	case "corporation":
		return GetCorporationName(contactID, esiclient)
	case "alliance":
		return GetAllianceName(contactID, esiclient)
	}
	return ""
}
