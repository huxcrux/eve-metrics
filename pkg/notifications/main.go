package notifications

import (
	"fmt"
	"os"

	"github.com/huxcrux/eve-metrics/pkg/config"
)

// CachedCollector collects metrics and serves them from cache
type NotificationController struct {
	Alliances    []NotificationControllerAlliances
	Corporations []NotificationControllerCoporation
	Characters   []NotificationControllerCharacter
	Webhook      string
}

type NotificationControllerAlliances struct {
	Character int32
	Token     string
	ID        int32
}

type NotificationControllerCoporation struct {
	Character int32
	Token     string
	ID        int32
}

type NotificationControllerCharacter struct {
	ID    int32
	Token string
}

// NewCachedCollector creates a new CachedCollector
func NewNotificationController(config config.Config) *NotificationController {

	characters := config.Notification
	localNotificationController := NotificationController{}

	// Webhook logic
	localNotificationController.Webhook = config.Discordwebhook

	// Alliance logic
	for _, alliance := range characters.Alliances {
		localAlliance := NotificationControllerAlliances{
			Character: alliance.Character,
			ID:        alliance.ID,
		}
		// Lookup the character token
		found := false
		for _, character := range characters.Characters {
			if character.ID == alliance.Character {
				found = true
				localAlliance.Token = character.Token
			}
		}
		if !found {
			fmt.Printf("Character with ID %d not found", alliance.Character)
			os.Exit(1)
		}
		localNotificationController.Alliances = append(localNotificationController.Alliances, localAlliance)
	}

	// Corporation logic
	for _, corporation := range characters.Corporations {
		localCorporation := NotificationControllerCoporation{
			Character: corporation.Character,
			ID:        corporation.ID,
		}
		// Lookup the character token
		found := false
		for _, character := range characters.Characters {
			if character.ID == corporation.Character {
				found = true
				localCorporation.Token = character.Token
			}
		}
		if !found {
			fmt.Printf("Character with ID %d not found", corporation.Character)
			os.Exit(1)
		}
		localNotificationController.Corporations = append(localNotificationController.Corporations, localCorporation)
	}

	// Character logic
	for _, character := range characters.Characters {
		localCharacter := NotificationControllerCharacter{
			ID:    character.ID,
			Token: character.Token,
		}
		localNotificationController.Characters = append(localNotificationController.Characters, localCharacter)
	}

	return &localNotificationController
}

// UpdateMetrics updates the cached metrics in the background
func (nc *NotificationController) Run() {

	// Check Alliance Contracts
	err := nc.FetchAllianceContacts()
	if err != nil {
		fmt.Println("Error fetching alliance contacts:", err)
	}

	// Check Alliance Corporations
	err = nc.FetchAllianceCorporations()
	if err != nil {
		fmt.Println("Error fetching alliance corporations:", err)
	}

	// Check Corporation Contacts
	err = nc.FetchCoporationContacts()
	if err != nil {
		fmt.Println("Error fetching corporation contacts:", err)
	}

	fmt.Println("Finished running notifications")
}
