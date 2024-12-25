package handlers

import (
	"encoding/json"
	"log"
	Metrics "myapp/metrics"
	"myapp/models"
	"myapp/utils"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RegisterImagesHandler struct {
	Metrics *Metrics.Metrics
	Imgs    *[]models.Image
}

func (rdh RegisterImagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getProcessingImages(w, r, rdh.Metrics, rdh.Imgs)
	case "POST":
		createImage(w, r, rdh.Metrics, rdh.Imgs)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getProcessingImages(w http.ResponseWriter, r *http.Request, m *Metrics.Metrics, imgs *[]models.Image) {
	now := time.Now()
	var processingImages []models.Image
	for _, img := range *imgs {
		if img.IMG_status == "Processing" {
			processingImages = append(processingImages, img)
		}
	}
	b, err := json.Marshal(processingImages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.Sleep_(200)
	m.Duration.With(prometheus.Labels{"method": "GET", "status": "200"}).Observe(time.Since(now).Seconds())
	log.Printf("Retrieved %d images being processed", len(processingImages))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func createImage(w http.ResponseWriter, r *http.Request, m *Metrics.Metrics, imgs *[]models.Image) {
	var img models.Image
	err := json.NewDecoder(r.Body).Decode(&img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*imgs = append(*imgs, img)
	ProcessingAndProcessedImagesCount := utils.GetProcessingAndProcessedImagesCount(*imgs)
	m.ProcessingImages.Set(float64(ProcessingAndProcessedImagesCount[0]))
	m.ProcessedImages.Set(float64(ProcessingAndProcessedImagesCount[1]))
	log.Printf("Image created: ID=%d, Format=%s, Resolution=%s Status=%s", img.ID, img.Format, img.Resolution, img.IMG_status)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created Image!"))
}
