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

	// start notification loop
	fmt.Println("Starting notification loop, will run every 60 seconds")
	for {
		notificationController.Run()
		time.Sleep(60 * time.Second)
	}

	// start reconciler loop
	// for {
	// 	// run the magic
	// 	cachedCollector.UpdateMetrics()
	// 	time.Sleep(60 * time.Second)
	// }
}
