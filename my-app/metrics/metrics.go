package Metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ProcessingImages prometheus.Gauge
	Info             *prometheus.GaugeVec
	ProcessedImages  prometheus.Gauge
	Duration         *prometheus.HistogramVec
	LoginDuration    prometheus.Summary
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		ProcessingImages: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "processing_images",
			Help:      "Number of images currently being processed.",
		}),
		Info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "info",
			Help:      "Information about the image processing environment.",
		}, []string{"version"}),
		ProcessedImages: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "processed_images",
			Help:      "Number of images processed.",
		}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "myapp",
			Name:      "request_duration_seconds",
			Help:      "Duration of the request for image processing.",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"status", "method"}),
		LoginDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  "myapp",
			Name:       "login_request_duration_seconds",
			Help:       "Duration of the login request.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
	reg.MustRegister(m.ProcessingImages, m.Info, m.ProcessedImages, m.Duration, m.LoginDuration)
	return m
}
