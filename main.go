package main

import (
	"log"
	"net/http"
	"time"

	"github.com/huxcrux/eve-metrics/pkg/collector"
	"github.com/huxcrux/eve-metrics/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	characters := config.ReadConfig()

	// Initialize the Metrics struct
	cachedCollector := collector.NewCachedCollector(characters)

	prometheus.MustRegister(cachedCollector)

	go func() {
		// Serve metrics via an HTTP server
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Serving metrics on :8080/metrics")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	// start reconciler loop
	for {
		// run the magic
		cachedCollector.UpdateMetrics()
		time.Sleep(60 * time.Second)
	}
}
