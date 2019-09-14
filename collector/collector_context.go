package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/routeros.v2"
	"mikrotik-exporter/config"
)

type collectorContext struct {
	ch     chan<- prometheus.Metric
	device *config.Device
	client *routeros.Client
}
