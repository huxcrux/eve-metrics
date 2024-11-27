package notifications

import (
	"fmt"
	"log"
	"time"

	"github.com/huxcrux/eve-metrics/pkg/config"
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

	for alliance := range nc.Alliances {

		// Make sure all values are set
		if nc.Alliances[alliance].Character == 0 || nc.Alliances[alliance].ID == 0 || nc.Alliances[alliance].Token == "" {
			fmt.Println("Alliance values not set")
			continue
		}

		// Create esiclient
		esiClient := esiClient.NewESIClient(nc.Alliances[alliance].Token)

		corporations, _, err := esiClient.Client.ESI.AllianceApi.GetAlliancesAllianceIdCorporations(esiClient.Ctx, int32(nc.Alliances[alliance].ID), nil)
		if err != nil {
			return err
		}

		allianceName := helpers.GetAllianceName(int32(nc.Alliances[alliance].ID), esiClient)

		// Notification logic
		AllianceMembersData := data.AllianceMember{
			AllianceID: int32(nc.Alliances[alliance].ID),
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
					AllianceID:      int32(nc.Alliances[alliance].ID),
					AllianceName:    allianceName,
					CorporationID:   corp,
					CorporationName: corpName,
					Event:           "added",
				}
				err = AllianceMemeberChange(change)
				if err != nil {
					fmt.Println("Error alerting on alliance member change:", err)
				}
			}
		}
		if len(removed) > 0 {
			changes = true
			for _, corp := range removed {
				corpName := helpers.GetCorporationName(corp, esiClient)
				fmt.Printf("Alliance: %s removed the corporation %s\n", allianceName, corpName)
				change := AllianceMemberChange{
					AllianceID:      int32(nc.Alliances[alliance].ID),
					AllianceName:    allianceName,
					CorporationID:   corp,
					CorporationName: corpName,
					Event:           "removed",
				}
				err = AllianceMemeberChange(change)
				if err != nil {
					fmt.Println("Error alerting on alliance member change:", err)
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

func AllianceMemeberChange(change AllianceMemberChange) error {

	// Get config
	configfile := config.ReadConfig()
	webhookURL := configfile.Discordwebhook

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
