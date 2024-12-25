package handlers

import (
	"log"
	Metrics "myapp/metrics"
	"myapp/models"
	"myapp/utils"

	"net/http"
	"strconv"
	"strings"
)

type ManageImagesHandler struct {
	Metrics *Metrics.Metrics
	Imgs    *[]models.Image
}

func (mdh ManageImagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		processImage(w, r, mdh.Metrics, mdh.Imgs)
	default:
		w.Header().Set("Allow", "PUT")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func processImage(w http.ResponseWriter, r *http.Request, m *Metrics.Metrics, imgs *[]models.Image) {
	path := strings.TrimPrefix(r.URL.Path, "/images/")
	id, err := strconv.Atoi(path)
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}
	for i := range *imgs {
		if (*imgs)[i].ID == id {
			(*imgs)[i].IMG_status = "Processed"
		}
	}
	ProcessingAndProcessedImagesCount := utils.GetProcessingAndProcessedImagesCount(*imgs)
	m.ProcessingImages.Set(float64(ProcessingAndProcessedImagesCount[0]))
	m.ProcessedImages.Set(float64(ProcessingAndProcessedImagesCount[1]))
	utils.Sleep_(1000)
	log.Printf("Image processed: ID=%d", id)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Processed Image.."))
}
