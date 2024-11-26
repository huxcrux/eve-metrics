package helpers

import esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"

func GetCharacterContactLabelName(lableID int64, charaterID int32, esiClient esiClient.ESIClient) string {
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

func GetAllianceContactLabelName(lableID int64, allianceID int32, esiClient esiClient.ESIClient) string {
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
