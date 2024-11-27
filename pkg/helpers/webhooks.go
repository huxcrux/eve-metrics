package helpers

import "github.com/huxcrux/eve-metrics/pkg/config"

func GetAllianceWebhooks(allianceID int32) []string {

	config := config.ReadConfig()

	// Generate webhook list
	webhooks := []string{}
	for _, item := range config.Webhooks {
		if item.AllAllainceSubscriptions {
			webhooks = append(webhooks, item.URL)
		} else {
			for _, alliance := range item.AllianceSubscriptions {
				if alliance == allianceID {
					webhooks = append(webhooks, item.URL)
				}
			}
		}
	}

	return webhooks
}
