package collector

import (
	"log"
	"sync"
	"time"

	"github.com/huxcrux/eve-metrics/pkg/helpers"
	"github.com/huxcrux/eve-metrics/pkg/models"
	"github.com/prometheus/client_golang/prometheus"
)

// CachedCollector collects metrics and serves them from cache
type CachedCollector struct {
	mu         sync.RWMutex
	cache      map[string]*prometheus.GaugeVec
	cacheTime  time.Time
	characters []models.Character
}

// NewCachedCollector creates a new CachedCollector
func NewCachedCollector(characters []models.CharacterInput) *CachedCollector {

	var localCharacters []models.Character
	for _, character := range characters {
		localCharacter, err := helpers.GetCharacterInfo(int32(character.ID), character.Token)
		if err == nil {
			//fmt.Printf("Character: %+v\n", localCharacter)
			localCharacters = append(localCharacters, localCharacter)
		} else {
			log.Println(err)
		}
	}

	return &CachedCollector{
		cache:      make(map[string]*prometheus.GaugeVec),
		cacheTime:  time.Now(),
		characters: localCharacters,
	}
}

// Describe sends metric descriptions (optional, leave empty for dynamic metrics)
func (cc *CachedCollector) Describe(ch chan<- *prometheus.Desc) {}

// Collect serves metrics from the cache
func (cc *CachedCollector) Collect(ch chan<- prometheus.Metric) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	// Serve metrics from the cache
	for _, gaugeVec := range cc.cache {
		gaugeVec.Collect(ch)
	}
}

// UpdateMetrics updates the cached metrics in the background
func (cc *CachedCollector) UpdateMetrics() {
	cc.mu.Lock()

	for character := range cc.characters {

		if character == 0 {
			// Update player count
			cc.FetchPlayersOnline(character)
		}

		// Update character info
		cc.FetchCharacter(character)

		// Update industry jobs
		cc.GenerateIndustryJobs(character)

		// Get contacts
		cc.FetchCharaterContacts(character)
		cc.FetchCorporationContacts(character)
		cc.FetchAllianceContacts(character)

		// Corporation
		cc.FetchAllianceCorporations(character)

		// Wallet
		// Gives a 403 for some reason
		// err := cc.FetchCharacterWallet()
		// if err != nil {
		// 	log.Println(err)
		// }

		// Jump fatigue
		err := cc.FetchJumpFatigue(character)
		if err != nil {
			log.Println("Error during FetchJumpFatigue: ", err)
		}

		// Jump Clone
		err = cc.FetchCharacterJumpClones(character)
		if err != nil {
			log.Println("error during FetchCharacterJumpClones: ", err)
		}

		// Live character Location and login info
		err = cc.FetchCharacterOnlineStatus(character)
		if err != nil {
			log.Println("Error during FetchCharacterOnlineStatus: ", err)
		}
		err = cc.FetchCharacterLocation(character)
		if err != nil {
			log.Println("Error during FetchCharacterLocation: ", err)
		}
		err = cc.FetchCharacterShip(character)
		if err != nil {
			log.Println("Error during FetchCharacterShip: ", err)
		}

		// Wars
		// TODO: fix a good way of filtering wars
		// err = cc.FetchWars()
		// if err != nil {
		// 	log.Println(err)
		// }

		// Loyalty Points
		err = cc.FetchCharacterLoyaltyPoints(character)
		if err != nil {
			log.Println("Error during FetchCharacterLoyaltyPoints: ", err)
		}

		// PI?
	}

	cc.cacheTime = time.Now()
	cc.mu.Unlock()
	log.Println("Metrics updated")
}
