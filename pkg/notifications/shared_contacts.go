package notifications

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/huxcrux/eve-metrics/pkg/config"
	"github.com/huxcrux/eve-metrics/pkg/discordwebhook"
	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/huxcrux/eve-metrics/pkg/models"
)

func ContactNotification(notification models.ContactNotificationInput) error {
	// Get config
	configfile := config.ReadConfig()
	webhookURL := configfile.Discordwebhook

	var title string
	var description string
	var color int
	if notification.Event == "added" {
		title = fmt.Sprintf("New %s contact", notification.EntityType)
		description = fmt.Sprintf(
			"The %s **%s** added a new contact an %s **%s** with standing **%.1f**.",
			notification.EntityType, notification.EntityName, notification.ContactType, notification.ContactName, *notification.Standing,
		)

		if len(notification.Labels) != 0 {
			description += fmt.Sprintf("\n\n**Labels:** %s", helpers.ListToCommaString(notification.Labels))
		}

		// Aqua color
		color = 0x1ABC9C
	} else if notification.Event == "removed" {
		title = fmt.Sprintf("%s contact removed", helpers.CapitalizeFirst(notification.EntityType))
		description = fmt.Sprintf(
			"The %s **%s** removed the contact %s: **%s**.",
			notification.EntityType, notification.EntityName, notification.ContactType, notification.ContactName,
		)

		if len(notification.Labels) != 0 {
			description += fmt.Sprintf("\n\n**Labels:** %s", helpers.ListToCommaString(notification.Labels))
		}

		// Orange color
		color = 0xE67E22

	} else {
		title = fmt.Sprintf("%s contact updated", helpers.CapitalizeFirst(notification.EntityType))
		description = fmt.Sprintf(
			"The %s **%s** updated the contact %s **%s** with standing **%.1f**.\n\n**Changes:**\n",
			notification.EntityType, notification.EntityName, notification.ContactType, notification.ContactName, *notification.Standing,
		)

		// Add standing change if applicable
		if notification.OldStanding != nil {
			// Compare with new value to make sure they are diffrent
			if *notification.OldStanding != *notification.Standing {
				description += fmt.Sprintf("• **Standing**: ~~%.1f~~ → **%.1f**\n", *notification.OldStanding, *notification.Standing)
			}
		}

		// Add label changes if applicable
		if notification.OldLabels != nil {
			// Compare with new value to make sure they are diffrent
			if !reflect.DeepEqual(*notification.OldLabels, notification.Labels) {
				oldLabels := helpers.ListToCommaString(*notification.OldLabels) // Helper to format labels as a readable string
				newLabels := helpers.ListToCommaString(notification.Labels)
				description += fmt.Sprintf("• **Labels**: ~~%s~~ → **%s**\n", oldLabels, newLabels)
			}
		}

		// Yellow color
		color = 0xF1C40F
	}
	// Create a webhook payload
	payload := discordwebhook.WebhookPayload{
		Content:  fmt.Sprintf("%s contact change", helpers.CapitalizeFirst(notification.EntityType)),
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
