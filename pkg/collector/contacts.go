package collector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) FetchCharaterContacts(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["player_contact"]; !exists {
		cc.cache["player_contact"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "player_contact",
				Help: "Player contact",
			},
			[]string{"character_name", "contact_name", "contact_id", "contact_type", "is_blocked", "is_watched", "labels"},
		)
	}

	// Fetch player count
	contact, _, err := cc.characters[characterID].ESIClient.Client.ESI.ContactsApi.GetCharactersCharacterIdContacts(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	characterName := helpers.GetCharacterName(int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient)
	for _, contact := range contact {
		contactName := ""
		switch contact.ContactType {
		case "character":
			contactName = helpers.GetCharacterName(contact.ContactId, cc.characters[characterID].ESIClient)
		case "corporation":
			contactName = helpers.GetCorporationName(contact.ContactId, cc.characters[characterID].ESIClient)
		case "alliance":
			contactName = helpers.GetAllianceName(contact.ContactId, cc.characters[characterID].ESIClient)
		}
		labelNames := []string{}
		if contact.LabelIds != nil {
			for _, labelId := range contact.LabelIds {
				labelNames = append(labelNames, helpers.GetCharacterContactLabelName(labelId, int32(cc.characters[characterID].ID), cc.characters[characterID].ESIClient))
			}
		}
		cc.cache["player_contact"].WithLabelValues(characterName, contactName, strconv.Itoa(int(contact.ContactId)), contact.ContactType, strconv.FormatBool(contact.IsBlocked), strconv.FormatBool(contact.IsWatched), strings.Join(labelNames, ",")).Set(float64(contact.Standing))

	}

	return nil
}

func (cc *CachedCollector) FetchCorporationContacts(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["corporation_contact"]; !exists {
		cc.cache["corporation_contact"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "corporation_contact",
				Help: "Corporation contact",
			},
			[]string{"contact_name", "contact_type", "is_watched", "label_ids"},
		)
	}

	// Fetch player count
	contact, _, err := cc.characters[characterID].ESIClient.Client.ESI.ContactsApi.GetCorporationsCorporationIdContacts(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].CorporationID), nil)
	if err != nil {
		return err
	}

	for _, contact := range contact {
		contactName := ""
		switch contact.ContactType {
		case "character":
			contactName = helpers.GetCharacterName(contact.ContactId, cc.characters[characterID].ESIClient)
		case "corporation":
			contactName = helpers.GetCorporationName(contact.ContactId, cc.characters[characterID].ESIClient)
		case "alliance":
			contactName = helpers.GetAllianceName(contact.ContactId, cc.characters[characterID].ESIClient)
		}
		cc.cache["corporation_contact"].WithLabelValues(contactName, contact.ContactType, strconv.FormatBool(contact.IsWatched), fmt.Sprint(contact.LabelIds)).Set(float64(contact.Standing))

	}

	return nil
}

func (cc *CachedCollector) FetchAllianceContacts(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["alliance_contact"]; !exists {
		cc.cache["alliance_contact"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "alliance_contact",
				Help: "Alliance contact",
			},
			[]string{"contact_name", "contact_type", "label_names"},
		)
	}

	// Fetch player count
	contact, _, err := cc.characters[characterID].ESIClient.Client.ESI.ContactsApi.GetAlliancesAllianceIdContacts(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].AllianceID), nil)
	if err != nil {
		return err
	}

	for _, contact := range contact {
		contactName := ""
		switch contact.ContactType {
		case "character":
			contactName = helpers.GetCharacterName(contact.ContactId, cc.characters[characterID].ESIClient)
		case "corporation":
			contactName = helpers.GetCorporationName(contact.ContactId, cc.characters[characterID].ESIClient)
		case "alliance":
			contactName = helpers.GetAllianceName(contact.ContactId, cc.characters[characterID].ESIClient)
		}
		labelNames := []string{}
		if contact.LabelIds != nil {
			for _, labelId := range contact.LabelIds {
				labelNames = append(labelNames, helpers.GetAllianceContactLabelName(labelId, int32(cc.characters[characterID].AllianceID), cc.characters[characterID].ESIClient))
			}
		}
		cc.cache["alliance_contact"].WithLabelValues(contactName, contact.ContactType, strings.Join(labelNames, ",")).Set(float64(contact.Standing))

	}

	return nil
}
