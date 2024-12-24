package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Image struct {
	ID         int    `json:"id"`
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
	IMG_status string `json:"img_status"`
}

type metrics struct {
	images          prometheus.Gauge
	info            *prometheus.GaugeVec
	processingCount *prometheus.CounterVec
	duration        *prometheus.HistogramVec
	loginDuration   prometheus.Summary
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		images: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "processing_images",
			Help:      "Number of images currently being processed.",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "info",
			Help:      "Information about the image processing environment.",
		},
			[]string{"version"}),
		processingCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "myapp",
			Name:      "total_images_processed",
			Help:      "Number of images processed.",
		}, []string{"type"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "myapp",
			Name:      "request_duration_seconds",
			Help:      "Duration of the request for image processing.",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"status", "method"}),
		loginDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  "myapp",
			Name:       "login_request_duration_seconds",
			Help:       "Duration of the login request.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
	reg.MustRegister(m.images, m.info, m.processingCount, m.duration, m.loginDuration)
	return m
}

var imgs []Image
var version string

func init() {
	version = "2.10.5"

	imgs = []Image{
		{1, "JPEG", "1920x1080", "Processing"},
		{2, "PNG", "1280x720", "Processing"},
	}
}

func main() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	m.images.Set(float64(len(imgs)))
	m.info.With(prometheus.Labels{"version": version}).Set(1)

	dMux := http.NewServeMux()
	rdh := registerImagesHandler{metrics: m}
	mdh := manageImagesHandler{metrics: m}

	lh := loginHandler{}
	mlh := middleware(lh, m)

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

type registerImagesHandler struct {
	metrics *metrics
}

func (rdh registerImagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getImages(w, r, rdh.metrics)
	case "POST":
		createImage(w, r, rdh.metrics)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getImages(w http.ResponseWriter, r *http.Request, m *metrics) {
	now := time.Now()

	b, err := json.Marshal(imgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sleep(200)

	m.duration.With(prometheus.Labels{"method": "GET", "status": "200"}).Observe(time.Since(now).Seconds())

	// Print the number of processed images to log for simulation
	log.Printf("Retrieved %d images", len(imgs))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func createImage(w http.ResponseWriter, r *http.Request, m *metrics) {
	var img Image

	err := json.NewDecoder(r.Body).Decode(&img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgs = append(imgs, img)

	// Update the processed images count
	m.images.Set(float64(len(imgs)))

	// Print out the new image details to simulate image processing
	log.Printf("Image created: ID=%d, Format=%s, Resolution=%s Status=%s", img.ID, img.Format, img.Resolution, img.IMG_status)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Processing Image!"))
}

func processImage(w http.ResponseWriter, r *http.Request, m *metrics) {
	path := strings.TrimPrefix(r.URL.Path, "/images/")

	id, err := strconv.Atoi(path)
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}

	var img Image
	err = json.NewDecoder(r.Body).Decode(&img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range imgs {
		if imgs[i].ID == id {
			// Simulate image processing
			imgs[i].IMG_status = "Processed"
		}
	}
	sleep(1000)

	// Increment the processing count
	m.processingCount.With(prometheus.Labels{"type": "resize"}).Inc()

	// Print the processed image details to the log
	log.Printf("Image processed: ID=%d, New Resolution=%s", img.ID, "Processed")

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Processing image..."))
}

type manageImagesHandler struct {
	metrics *metrics
}

func (mdh manageImagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		processImage(w, r, mdh.metrics)
	default:
		w.Header().Set("Allow", "PUT")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func sleep(ms int) {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	n := rand.Intn(ms + now.Second())
	time.Sleep(time.Duration(n) * time.Millisecond)
}

type loginHandler struct{}

func (l loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sleep(200)
	w.Write([]byte("Welcome to the image processing app!"))
}

func middleware(next http.Handler, m *metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		m.loginDuration.Observe(time.Since(now).Seconds())
	})
}
