package helpers

import (
	"fmt"
	"log"

	"github.com/antihax/goesi/esi"
	"github.com/huxcrux/eve-metrics/pkg/config"
	"github.com/huxcrux/eve-metrics/pkg/data"
	"github.com/huxcrux/eve-metrics/pkg/discordwebhook"
)

type CompletedIndustryJobs struct {
	Job             esi.GetCharactersCharacterIdIndustryJobs200Ok
	BlueprintName   string
	StructureName   string
	StructureSystem string
	StructureOwner  string
	CharacterName   string
	Activity        string
}

func IndustryJobCompleted(indujob CompletedIndustryJobs) error {

	// Check if job exists in data file
	datafile := data.ReadData()

	found := false

	for _, job := range datafile.CompletedIndustryJobs {
		if job == indujob.Job.JobId {
			found = true
			break
		}
	}

	if !found {
		datafile.CompletedIndustryJobs = append(datafile.CompletedIndustryJobs, indujob.Job.JobId)
		data.WriteData(datafile)

		// Send alert
		// Get config
		configfile := config.ReadConfig()
		webhookURL := configfile.Discordwebhook

		// Calculate success percentage
		successPercentage := float64(indujob.Job.SuccessfulRuns) / float64(indujob.Job.Runs) * 100

		// Create a webhook payload
		payload := discordwebhook.WebhookPayload{
			Content:  "Industry Job Completed",
			Username: "EVE Bot",
			Embeds: []discordwebhook.Embed{
				{
					Title: "Industry Job Completed",
					Description: fmt.Sprintf(
						"%s completed an **%s** activity using the blueprint **%s** in the structure **%s** located in **%s**, owned by **%s**.",
						indujob.CharacterName, indujob.Activity, indujob.BlueprintName, indujob.StructureName, indujob.StructureSystem, indujob.StructureOwner,
					),
					Color: 0x00FF00, // Green color
					Fields: []discordwebhook.Field{
						{Name: "Character", Value: indujob.CharacterName, Inline: true},
						{Name: "Activity", Value: indujob.Activity, Inline: true},
						{Name: "Blueprint", Value: indujob.BlueprintName, Inline: true},
						{Name: "Structure", Value: indujob.StructureName, Inline: true},
						{Name: "System", Value: indujob.StructureSystem, Inline: true},
						{Name: "Owner", Value: indujob.StructureOwner, Inline: true},
						{Name: "Job ID", Value: fmt.Sprintf("%d", indujob.Job.JobId), Inline: true},
						{Name: "Runs", Value: fmt.Sprintf("%d (%.2f%% successful)", indujob.Job.Runs, successPercentage), Inline: true},
						{Name: "Cost", Value: fmt.Sprintf("%.2f ISK", indujob.Job.Cost), Inline: true},
						{Name: "Completed", Value: indujob.Job.EndDate.Format("2006-01-02 15:04:05"), Inline: false},
					},
				},
			},
		}

		// Send the webhook
		err := discordwebhook.SendWebhook(webhookURL, payload)
		if err != nil {
			log.Fatalf("Failed to send webhook: %v", err)
		}

		log.Println("Webhook sent successfully!")
	}

	return nil
}
