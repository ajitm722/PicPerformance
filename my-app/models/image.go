package models

type Image struct {
	ID         int    `json:"id"`
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
	IMG_status string `json:"img_status"`
}

