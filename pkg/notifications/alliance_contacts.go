package notifications

import (
	"fmt"
	"reflect"

	"github.com/huxcrux/eve-metrics/pkg/data"
	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/huxcrux/eve-metrics/pkg/models"
)

type AllianceContactChanges struct {
	Added   []data.AllianceContact
	Removed []data.AllianceContact
	Changed []AllianceContactChange
}

type AllianceContactChange struct {
	Old data.AllianceContact
	New data.AllianceContact
}

func (nc *NotificationController) FetchAllianceContacts() error {

	for alliance := range nc.Alliances {
		// Make sure all values are set
		if nc.Alliances[alliance].Character == 0 || nc.Alliances[alliance].ID == 0 || nc.Alliances[alliance].Token == "" {
			fmt.Println("Alliance values not set")
			continue
		}

		// Create esiclient
		esiClient := esiClient.NewESIClient(nc.Alliances[alliance].Token)

		// Fetch alliance contacts
		contacts, _, err := esiClient.Client.ESI.ContactsApi.GetAlliancesAllianceIdContacts(esiClient.Ctx, int32(nc.Alliances[alliance].ID), nil)
		if err != nil {
			return err
		}

		// Notification logic
		allianceContacts := []data.AllianceContact{}
		for _, contact := range contacts {
			allianceContacts = append(allianceContacts, data.AllianceContact{
				ContactId:   contact.ContactId,
				ContactType: contact.ContactType,
				Standing:    contact.Standing,
				LabelIds:    contact.LabelIds,
			})
		}
		allianceContact := data.AllianceContacts{
			AllianceID: nc.Alliances[alliance].ID,
			Contacts:   allianceContacts,
		}

		allianceName := helpers.GetAllianceName(int32(nc.Alliances[alliance].ID), esiClient)
		changes := CompareAllianceContacts(allianceContact)
		ChangeFound := false
		if len(changes.Added) > 0 {
			ChangeFound = true
			for _, contact := range changes.Added {
				contactName := helpers.GetContactName(contact.ContactType, contact.ContactId, esiClient)
				labels := []string{}
				for _, label := range contact.LabelIds {
					locallabel := helpers.GetAllianceContactLabelName(label, nc.Alliances[alliance].ID, esiClient)
					labels = append(labels, locallabel)
				}
				fmt.Printf("Alliance: **%s** added the %s **%s** as a contact\n", allianceName, contact.ContactType, contactName)
				an := models.ContactNotificationInput{
					Event:       "added",
					EntityType:  "alliance",
					EntityName:  allianceName,
					ContactType: contact.ContactType,
					ContactName: contactName,
					Standing:    &contact.Standing,
					OldStanding: nil,
					Labels:      labels,
					OldLabels:   nil,
				}
				err = ContactNotification(an)
				if err != nil {
					fmt.Println("Error alerting on alliance contact change:", err)
				}
			}
		}
		if len(changes.Removed) > 0 {
			ChangeFound = true
			for _, contact := range changes.Removed {
				contactName := helpers.GetContactName(contact.ContactType, contact.ContactId, esiClient)
				labels := []string{}
				for _, label := range contact.LabelIds {
					locallabel := helpers.GetAllianceContactLabelName(label, nc.Alliances[alliance].ID, esiClient)
					labels = append(labels, locallabel)
				}
				fmt.Printf("Alliance: **%s** removed the %s **%s** as a contact\n", allianceName, contact.ContactType, contactName)
				an := models.ContactNotificationInput{
					Event:       "removed",
					EntityType:  "alliance",
					EntityName:  allianceName,
					ContactType: contact.ContactType,
					ContactName: contactName,
					Standing:    nil,
					OldStanding: nil,
					Labels:      labels,
					OldLabels:   nil,
				}
				err = ContactNotification(an)
				if err != nil {
					fmt.Println("Error alerting on alliance contact change:", err)
				}
			}
		}
		if len(changes.Changed) > 0 {
			ChangeFound = true
			for _, contact := range changes.Changed {
				contactName := helpers.GetContactName(contact.New.ContactType, contact.New.ContactId, esiClient)
				labels := []string{}
				for _, label := range contact.New.LabelIds {
					locallabel := helpers.GetAllianceContactLabelName(label, nc.Alliances[alliance].ID, esiClient)
					labels = append(labels, locallabel)
				}
				oldLabels := []string{}
				for _, label := range contact.Old.LabelIds {
					locallabel := helpers.GetAllianceContactLabelName(label, nc.Alliances[alliance].ID, esiClient)
					oldLabels = append(oldLabels, locallabel)
				}
				fmt.Printf("Alliance: **%s** Updated the %s **%s** as a contact\n", allianceName, contact.New.ContactType, contactName)
				an := models.ContactNotificationInput{
					Event:       "updated",
					EntityType:  "alliance",
					EntityName:  allianceName,
					ContactType: contact.New.ContactType,
					ContactName: contactName,
					Standing:    &contact.New.Standing,
					OldStanding: &contact.Old.Standing,
					Labels:      labels,
					OldLabels:   &oldLabels,
				}
				err = ContactNotification(an)
				if err != nil {
					fmt.Println("Error alerting on alliance contact change:", err)
				}
			}
		}

		if ChangeFound {
			fmt.Println("Changes found in alliance contacts")
			datafile := data.ReadData()
			found := false
			for i, alliance := range datafile.AllianceContacts {
				if alliance.AllianceID == allianceContact.AllianceID {
					datafile.AllianceContacts[i] = allianceContact
					found = true
					break
				}
			}
			if !found {
				datafile.AllianceContacts = append(datafile.AllianceContacts, allianceContact)
			}
			data.WriteData(datafile)

		}
	}

	return nil
}

func CompareAllianceContacts(new data.AllianceContacts) AllianceContactChanges {

	// Read data file
	var old data.AllianceContacts
	olddata := data.ReadData()
	for _, alliance := range olddata.AllianceContacts {
		if alliance.AllianceID == new.AllianceID {
			old = alliance
			break
		}
	}

	changes := AllianceContactChanges{}

	// Create maps for quick lookup
	oldMap := make(map[int32]data.AllianceContact)
	newMap := make(map[int32]data.AllianceContact)

	for _, contact := range old.Contacts {
		oldMap[contact.ContactId] = contact
	}
	for _, contact := range new.Contacts {
		newMap[contact.ContactId] = contact
	}

	// Compare old to new
	for id, oldContact := range oldMap {
		newContact, exists := newMap[id]
		if !exists {
			// Removed contact
			changes.Removed = append(changes.Removed, oldContact)
		} else if !reflect.DeepEqual(oldContact, newContact) {
			// Changed contact
			changes.Changed = append(changes.Changed, AllianceContactChange{Old: oldContact, New: newContact})
		}
	}

	// Find added contacts
	for id, newContact := range newMap {
		if _, exists := oldMap[id]; !exists {
			changes.Added = append(changes.Added, newContact)
		}
	}

	return changes
}
