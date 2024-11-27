package notifications

import (
	"fmt"
	"log"
	"time"

	"github.com/huxcrux/eve-metrics/pkg/data"
	"github.com/huxcrux/eve-metrics/pkg/discordwebhook"
	esiClient "github.com/huxcrux/eve-metrics/pkg/esi_client"
	"github.com/huxcrux/eve-metrics/pkg/helpers"
)

type AllianceMemberChange struct {
	AllianceID      int32
	AllianceName    string
	CorporationID   int32
	CorporationName string
	Event           string
}

func (nc *NotificationController) FetchAllianceCorporations() error {

	// Get the first alliance token
	// This is just to auth the cacheing proxy
	if len(nc.Alliances) == 0 {
		return fmt.Errorf("no alliances found")
	}
	esiClient := esiClient.NewESIClient(nc.Alliances[0].Token)

	// Get all alliances
	alliances, _, err := esiClient.Client.ESI.AllianceApi.GetAlliances(esiClient.Ctx, nil)
	if err != nil {
		fmt.Println("Error fetching alliances:", err)
	}

	fmt.Println("Number of alliances found: ", len(alliances))

	addedAliances, removedAliances := CompareAlliances(alliances)

	// Notification logic
	if len(addedAliances) > 0 {
		for _, alliance := range addedAliances {
			allianceName := helpers.GetAllianceName(alliance, esiClient)
			fmt.Printf("Alliance: %s added\n", allianceName)
			change := AllianceMemberChange{
				AllianceID:      alliance,
				AllianceName:    allianceName,
				CorporationID:   0,
				CorporationName: "",
				Event:           "added",
			}
			// Get webhook list
			webhooks := helpers.GetAllianceWebhooks(alliance)
			// Send all webhooks
			for _, webhookitem := range webhooks {
				err = AllianceChange(change, webhookitem)
				if err != nil {
					fmt.Println("Error alerting on alliance member change:", err)
				}
			}
		}
	}

	if len(removedAliances) > 0 {
		for _, alliance := range removedAliances {
			allianceName := helpers.GetAllianceName(alliance, esiClient)
			fmt.Printf("Alliance: %s removed\n", allianceName)
			change := AllianceMemberChange{
				AllianceID:      alliance,
				AllianceName:    allianceName,
				CorporationID:   0,
				CorporationName: "",
				Event:           "removed",
			}
			// Get webhook list
			webhooks := helpers.GetAllianceWebhooks(alliance)
			// Send all webhooks
			for _, webhook := range webhooks {
				err = AllianceChange(change, webhook)
				if err != nil {
					fmt.Println("Error alerting on alliance member change:", err)
				}
			}
		}
	}

	// Check for new or removed alliances

	for _, alliance := range alliances {

		// Make sure all values are set
		corporations, _, err := esiClient.Client.ESI.AllianceApi.GetAlliancesAllianceIdCorporations(esiClient.Ctx, int32(alliance), nil)
		if err != nil {
			return err
		}

		webhooks := helpers.GetAllianceWebhooks(int32(alliance))

		allianceName := helpers.GetAllianceName(int32(alliance), esiClient)

		// Notification logic
		AllianceMembersData := data.AllianceMember{
			AllianceID: int32(alliance),
			Members:    corporations,
		}

		added, removed := CompareAllianceMembers(AllianceMembersData)
		changes := false
		if len(added) > 0 {
			changes = true
			for _, corp := range added {
				corpName := helpers.GetCorporationName(corp, esiClient)
				fmt.Printf("Alliance: %s added the corporation %s\n", allianceName, corpName)
				change := AllianceMemberChange{
					AllianceID:      int32(alliance),
					AllianceName:    allianceName,
					CorporationID:   corp,
					CorporationName: corpName,
					Event:           "added",
				}
				for _, webhook := range webhooks {
					err = AllianceMemeberChange(change, webhook)
					if err != nil {
						fmt.Println("Error alerting on alliance member change:", err)
					}
				}
			}
		}
		if len(removed) > 0 {
			changes = true
			for _, corp := range removed {
				corpName := helpers.GetCorporationName(corp, esiClient)
				fmt.Printf("Alliance: %s removed the corporation %s\n", allianceName, corpName)
				change := AllianceMemberChange{
					AllianceID:      int32(alliance),
					AllianceName:    allianceName,
					CorporationID:   corp,
					CorporationName: corpName,
					Event:           "removed",
				}
				for _, webhook := range webhooks {
					err = AllianceMemeberChange(change, webhook)
					if err != nil {
						fmt.Println("Error alerting on alliance member change:", err)
					}
				}
			}
		}

		if changes {
			datafile := data.ReadData()
			found := false
			for i, alliance := range datafile.AllianceMembers {
				if alliance.AllianceID == AllianceMembersData.AllianceID {
					datafile.AllianceMembers[i] = AllianceMembersData
					found = true
					break
				}
			}
			if !found {
				datafile.AllianceMembers = append(datafile.AllianceMembers, AllianceMembersData)
			}
			data.WriteData(datafile)
		}
	}

	return nil
}

// CompareAlliances compares two AllianceMember structs and returns the added and removed members.
func CompareAlliances(updated []int32) (added, removed []int32) {

	saved := data.ReadData()
	var old []int32
	for _, alliance := range saved.AllianceMembers {
		old = append(old, alliance.AllianceID)
	}

	oldMap := make(map[int32]bool)
	newMap := make(map[int32]bool)

	// Populate maps for quick lookup
	for _, alliance := range old {
		oldMap[alliance] = true
	}
	for _, alliance := range updated {
		newMap[alliance] = true
	}

	// Check for added members
	for _, alliance := range updated {
		if !oldMap[alliance] {
			added = append(added, alliance)
		}
	}

	// Check for removed members
	for _, alliance := range old {
		if !newMap[alliance] {
			removed = append(removed, alliance)
		}
	}

	return added, removed
}

// CompareAllianceMembers compares two AllianceMember structs and returns the added and removed members.
func CompareAllianceMembers(updated data.AllianceMember) (added, removed []int32) {

	svaed := data.ReadData()
	var old data.AllianceMember
	for _, alliance := range svaed.AllianceMembers {
		if alliance.AllianceID == updated.AllianceID {
			old = alliance
			break
		}
	}

	oldMap := make(map[int32]bool)
	newMap := make(map[int32]bool)

	// Populate maps for quick lookup
	for _, member := range old.Members {
		oldMap[member] = true
	}
	for _, member := range updated.Members {
		newMap[member] = true
	}

	// Check for added members
	for _, member := range updated.Members {
		if !oldMap[member] {
			added = append(added, member)
		}
	}

	// Check for removed members
	for _, member := range old.Members {
		if !newMap[member] {
			removed = append(removed, member)
		}
	}

	return added, removed
}

func AllianceMemeberChange(change AllianceMemberChange, webhookURL string) error {

	var title string
	var description string
	var color int
	if change.Event == "added" {
		change.Event = "joined"
		title = "Alliance Member Added"
		description = fmt.Sprintf(
			"The corporation **%s** joined **%s**.",
			change.CorporationName, change.AllianceName,
		)
		// Aqua color
		color = 0x1ABC9C
	} else {
		change.Event = "left"
		title = "Alliance Member Removed"
		description = fmt.Sprintf(
			"The corporation **%s** left **%s**.",
			change.CorporationName, change.AllianceName,
		)
		// Orange color
		color = 0xE67E22

	}
	// Create a webhook payload
	payload := discordwebhook.WebhookPayload{
		Content:  "Alliance Member Change",
		Username: "EVE Bot",
		Embeds: []discordwebhook.Embed{
			{
				Title:       title,
				Description: description,
				Color:       color,
			},
		},
	}

	// Send the webhook, in case of failure retry every 10 seconds, this is blocking, maybe sping it off to a go routine?
	for {
		err := discordwebhook.SendWebhook(webhookURL, payload)
		if err != nil {
			log.Printf("Failed to send webhook: %v", err)
		} else {
			log.Println("Webhook sent successfully!")
			break
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}

func AllianceChange(change AllianceMemberChange, webhookURL string) error {

	var title string
	var description string
	var color int
	if change.Event == "added" {
		title = "Alliance Created"
		description = fmt.Sprintf(
			"The alliance **%s** was created.",
			change.AllianceName,
		)
		// Aqua color
		color = 0x1ABC9C
	} else {
		title = "Alliance Disbanded"
		description = fmt.Sprintf(
			"The alliance **%s** is now disbanded.",
			change.AllianceName,
		)
		// Orange color
		color = 0xE67E22

	}
	// Create a webhook payload
	payload := discordwebhook.WebhookPayload{
		Content:  "Alliance Change",
		Username: "EVE Bot",
		Embeds: []discordwebhook.Embed{
			{
				Title:       title,
				Description: description,
				Color:       color,
			},
		},
	}

	// Send the webhook, in case of failure retry every 10 seconds, this is blocking, maybe sping it off to a go routine?
	for {
		err := discordwebhook.SendWebhook(webhookURL, payload)
		if err != nil {
			log.Printf("Failed to send webhook: %v", err)
		} else {
			log.Println("Webhook sent successfully!")
			break
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}
