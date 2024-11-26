package collector

import (
	"fmt"

	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func (cc *CachedCollector) GenerateIndustryJobs(characterID int) error {
	// Update metric values dynamically
	if _, exists := cc.cache["industry_job"]; !exists {
		cc.cache["industry_job"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "industry_job",
				Help: "Active industry jobs",
			},
			[]string{
				"activity_name",
				"blueprint_name",
				"station_name",
				"solar_system",
				"owner",
				"cost",
				"duration",
				"character",
				"job_id",
				"licensed_runs",
				"probability",
				"runs",
				"status",
				"successful_runs",
			},
		)
	}

	// fetch industry jobs
	induJobs, _, err := cc.characters[characterID].ESIClient.Client.ESI.IndustryApi.GetCharactersCharacterIdIndustryJobs(cc.characters[characterID].ESIClient.Ctx, int32(cc.characters[characterID].ID), nil)
	if err != nil {
		return err
	}

	for job := range induJobs {
		// Calculate seconds left
		time := induJobs[job].EndDate.Unix()
		blueprintName := helpers.GetBlueprintName(int(induJobs[job].BlueprintTypeId), cc.characters[characterID].ESIClient)
		structureName := helpers.GetStructureName(induJobs[job].StationId, cc.characters[characterID].ESIClient)
		structureSystem := helpers.GetStructureSystem(induJobs[job].StationId, cc.characters[characterID].ESIClient)
		structureOwner := helpers.GetStructureOwner(induJobs[job].StationId, cc.characters[characterID].ESIClient)
		characterName := helpers.GetCharacterName(induJobs[job].InstallerId, cc.characters[characterID].ESIClient)
		activity := helpers.GetIndustryActivityName(induJobs[job].ActivityId)
		cc.cache["industry_job"].WithLabelValues(
			activity,
			blueprintName,
			structureName,
			structureSystem,
			structureOwner,
			fmt.Sprintf("%f", induJobs[job].Cost),
			fmt.Sprintf("%d", induJobs[job].Duration),
			characterName,
			fmt.Sprintf("%d", induJobs[job].JobId),
			fmt.Sprintf("%d", induJobs[job].LicensedRuns),
			fmt.Sprintf("%f", induJobs[job].Probability),
			fmt.Sprintf("%d", induJobs[job].Runs),
			induJobs[job].Status,
			fmt.Sprintf("%d", induJobs[job].SuccessfulRuns),
		).Set(float64(time))
		// fmt.Printf("Industry Job: %+v\n", induJobs[job])
		// fmt.Printf("Structure Name: %s, System: %s, Owner: %s\n", structureName, structureSystem, structureOwner)
	}

	return nil
}
