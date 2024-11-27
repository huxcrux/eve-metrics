package notifications

import (
	"fmt"

	"github.com/go-test/deep"
	"github.com/huxcrux/eve-metrics/pkg/data"
	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/huxcrux/eve-metrics/pkg/models"
)

type CoporationContactChanges struct {
	Added   []data.CorporationContact
	Removed []data.CorporationContact
	Changed []CoporationContactChange
}

type CoporationContactChange struct {
	Old data.CorporationContact
	New data.CorporationContact
}

func (nc *NotificationController) FetchCoporationContacts() error {

	for corporation := range nc.Corporations {

		// Make sure all values are set
		if nc.Corporations[corporation].Character == 0 || nc.Corporations[corporation].ID == 0 || nc.Corporations[corporation].Token == "" {
			fmt.Println("Corporation values not set")
			continue
		}

		// Create esiclient
		esiClient := esiClient.NewESIClient(nc.Corporations[corporation].Token)

		// Fetch corporation contacts
		contacts, _, err := esiClient.Client.ESI.ContactsApi.GetCorporationsCorporationIdContacts(esiClient.Ctx, int32(nc.Corporations[corporation].ID), nil)
		if err != nil {
			return err
		}

		// Notification logic

		localContacts := []data.CorporationContact{}
		for _, contact := range contacts {
			localContacts = append(localContacts, data.CorporationContact{
				ID:          contact.ContactId,
				Standing:    contact.Standing,
				ContactType: contact.ContactType,
				LabelIds:    contact.LabelIds})
		}
		corporationContacts := data.CorporationContacts{
			CorporationID: nc.Corporations[corporation].ID,
			Contacts:      localContacts,
		}

		copoationName := helpers.GetCorporationName(int32(nc.Corporations[corporation].ID), esiClient)
		// Compare corporation contacts
		changes := CompareCorporationContacts(corporationContacts)
		ChangeFound := false
		if len(changes.Added) > 0 {
			ChangeFound = true
			for _, contact := range changes.Added {
				contactName := helpers.GetContactName(contact.ContactType, contact.ID, esiClient)
				fmt.Printf("Corporation: %s added the contact %s\n", copoationName, contactName)
				// Send notification
				labels := []string{}
				if contact.LabelIds != nil {
					for _, labelId := range contact.LabelIds {
						labels = append(labels, helpers.GetCorporationContactLabelName(labelId, int32(nc.Corporations[corporation].ID), esiClient))
					}
				}
				ContactNotification(models.ContactNotificationInput{
					Event:       "added",
					EntityType:  "corporation",
					EntityName:  copoationName,
					ContactType: contact.ContactType,
					ContactName: contactName,
					Standing:    &contact.Standing,
					Labels:      labels,
				})
			}
		}

		if len(changes.Removed) > 0 {
			ChangeFound = true
			for _, contact := range changes.Removed {
				contactName := helpers.GetContactName(contact.ContactType, contact.ID, esiClient)
				fmt.Printf("Corporation: %s removed the contact %s\n", copoationName, contactName)
				// Send notification
				labels := []string{}
				if contact.LabelIds != nil {
					for _, labelId := range contact.LabelIds {
						labels = append(labels, helpers.GetCorporationContactLabelName(labelId, int32(nc.Corporations[corporation].ID), esiClient))
					}
				}
				ContactNotification(models.ContactNotificationInput{
					Event:       "removed",
					EntityType:  "corporation",
					EntityName:  copoationName,
					ContactType: contact.ContactType,
					ContactName: contactName,
					Labels:      labels,
				})
			}
		}

		if len(changes.Changed) > 0 {
			ChangeFound = true
			for _, contact := range changes.Changed {
				contactName := helpers.GetContactName(contact.New.ContactType, contact.New.ID, esiClient)
				fmt.Printf("Corporation: %s changed the contact %s\n", copoationName, contactName)
				// Send notification
				labels := []string{}
				if contact.New.LabelIds != nil {
					for _, labelId := range contact.New.LabelIds {
						labels = append(labels, helpers.GetCorporationContactLabelName(labelId, int32(nc.Corporations[corporation].ID), esiClient))
					}
				}
				ContactNotification(models.ContactNotificationInput{
					Event:       "updated",
					EntityType:  "corporation",
					EntityName:  copoationName,
					ContactType: contact.New.ContactType,
					ContactName: contactName,
					Standing:    &contact.New.Standing,
					OldStanding: &contact.Old.Standing,
					Labels:      labels,
					OldLabels:   nil,
				})
			}
		}

		if ChangeFound {
			// Save changes
			datafile := data.ReadData()
			found := false
			for i, corporation := range datafile.CorporationContacts {
				if corporation.CorporationID == corporationContacts.CorporationID {
					datafile.CorporationContacts[i] = corporationContacts
					found = true
					break
				}
			}
			if !found {
				datafile.CorporationContacts = append(datafile.CorporationContacts, corporationContacts)
			}
			data.WriteData(datafile)
		}
	}

	return nil
}

func CompareCorporationContacts(new data.CorporationContacts) CoporationContactChanges {

	// Read data file
	var old data.CorporationContacts
	olddata := data.ReadData()
	for _, corporation := range olddata.CorporationContacts {
		if corporation.CorporationID == new.CorporationID {
			old = corporation
			break
		}
	}

	changes := CoporationContactChanges{}

	// Create maps for quick lookup
	oldMap := make(map[int32]data.CorporationContact)
	newMap := make(map[int32]data.CorporationContact)

	for _, contact := range old.Contacts {
		oldMap[contact.ID] = contact
	}
	for _, contact := range new.Contacts {
		newMap[contact.ID] = contact
	}

	// Compare old to new
	for id, oldContact := range oldMap {
		newContact, exists := newMap[id]
		if !exists {
			// Removed contact
			changes.Removed = append(changes.Removed, oldContact)
		}
		diff := deep.Equal(oldContact, newContact)
		if diff != nil {
			// Changed contact
			fmt.Println(diff)
			changes.Changed = append(changes.Changed, CoporationContactChange{Old: oldContact, New: newContact})
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
