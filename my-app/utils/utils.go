package utils

import (
	"math/rand"
	"myapp/models"
	"time"
)

func Sleep_(ms int) {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	n := rand.Intn(ms + now.Second())
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func GetProcessingAndProcessedImagesCount(imgs []models.Image) [2]int {
	processing, processed := 0, 0
	for _, img := range imgs {
		if img.IMG_status == "Processing" {
			processing++
		} else if img.IMG_status == "Processed" {
			processed++
		}
	}
	return [2]int{processing, processed}
}
