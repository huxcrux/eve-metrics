package main

import (
	"fmt"
	"time"

	"github.com/huxcrux/eve-metrics/pkg/config"
	"github.com/huxcrux/eve-metrics/pkg/notifications"
)

func main() {

	fmt.Println("Starting eve-metrics")

	// Read the configuration file
	config := config.ReadConfig()

	// Initialize the Metrics struct
	//cachedCollector := collector.NewCachedCollector(config.Characters)
	notificationController := notifications.NewNotificationController(config)

	//prometheus.MustRegister(cachedCollector)

	// go func() {
	// 	// Serve metrics via an HTTP server
	// 	http.Handle("/metrics", promhttp.Handler())
	// 	log.Println("Serving metrics on :8080/metrics")
	// 	if err := http.ListenAndServe(":8080", nil); err != nil {
	// 		log.Fatalf("Error starting HTTP server: %v", err)
	// 	}
	// }()

	// Start all other notification loops
	fmt.Println("Starting Alliance/Coporation contact loop")
	counter := 0
	for {
		fmt.Println("Running alliance/Coporation contact loop")
		notificationController.Run()
		time.Sleep(60 * time.Second)
		// TODO: replace by a timestamp stored in data file
		if counter > 15 {
			fmt.Println("Starting Alliance Members loop")
			notificationController.RunAllianceMembers()
			counter = 0
		}
		counter++
	}

	// start reconciler loop
	// for {
	// 	// run the magic
	// 	cachedCollector.UpdateMetrics()
	// 	time.Sleep(60 * time.Second)
	// }
}
