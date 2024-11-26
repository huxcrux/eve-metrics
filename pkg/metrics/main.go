package metrics

import (
	"fmt"
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics struct to hold dynamic metrics
type Metrics struct {
	mutex    sync.Mutex
	Counters map[string]*prometheus.CounterVec
	Gauges   map[string]*prometheus.GaugeVec
}

// NewMetrics initializes the custom metrics struct
func NewMetrics() *Metrics {
	return &Metrics{
		Counters: make(map[string]*prometheus.CounterVec),
		Gauges:   make(map[string]*prometheus.GaugeVec),
	}
}

// AddCounter adds a new counter metric dynamically
func (m *Metrics) AddCounter(name, help string, labels []string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.Counters[name]; exists {
		log.Printf("Counter %s already exists", name)
		return
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	}, labels)

	prometheus.MustRegister(counter)
	m.Counters[name] = counter
}

// CreateGauge creates a new GaugeVec with the specified name, help text, and labels
func (m *Metrics) CreateGauge(name string, help string, labels []string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.Gauges[name]; exists {
		return fmt.Errorf("gauge with name '%s' already exists", name)
	}

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, labels)

	err := prometheus.Register(gauge)
	if err != nil {
		return fmt.Errorf("failed to register gauge: %v", err)
	}

	m.Gauges[name] = gauge
	return nil
}

// UpdateGauge updates the value of a gauge by setting it to the given value
func (m *Metrics) UpdateGauge(name string, labels prometheus.Labels, value float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	gauge, exists := m.Gauges[name]
	if !exists {
		return fmt.Errorf("gauge with name '%s' does not exist", name)
	}

	gauge.With(labels).Set(value)
	return nil
}

// DeleteGauge removes a gauge by name
func (m *Metrics) DeleteGauge(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	gauge, exists := m.Gauges[name]
	if !exists {
		return fmt.Errorf("gauge with name '%s' does not exist", name)
	}

	prometheus.Unregister(gauge)
	delete(m.Gauges, name)
	return nil
}
