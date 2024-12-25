package main

import (
	"log"

	"myapp/handlers"
	Metrics "myapp/metrics"
	"myapp/models"
	"myapp/utils"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var imgs []models.Image
var version string

func init() {
	version = "2.10.5"
}

func main() {
	reg := prometheus.NewRegistry()
	m := Metrics.NewMetrics(reg)

	ProcessingAndProcessedImagesCount := utils.GetProcessingAndProcessedImagesCount(imgs)
	m.ProcessingImages.Set(float64(ProcessingAndProcessedImagesCount[0]))
	m.ProcessedImages.Set(float64(ProcessingAndProcessedImagesCount[1]))
	m.Info.With(prometheus.Labels{"version": version}).Set(1)

	dMux := http.NewServeMux()
	rdh := handlers.RegisterImagesHandler{Metrics: m, Imgs: &imgs}
	mdh := handlers.ManageImagesHandler{Metrics: m, Imgs: &imgs}
	lh := handlers.LoginHandler{}
	mlh := handlers.Middleware(lh, m)

	dMux.Handle("/images", rdh)
	dMux.Handle("/images/", mdh)
	dMux.Handle("/login", mlh)

	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", dMux))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8081", pMux))
	}()

	select {}
}
