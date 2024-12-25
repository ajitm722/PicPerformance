package handlers

import (
	Metrics "myapp/metrics"
	"myapp/utils"
	"net/http"
	"time"
)

type LoginHandler struct{}

func (l LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.Sleep_(200)
	w.Write([]byte("Welcome to the image processing app!"))
}

func Middleware(next http.Handler, m *Metrics.Metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		m.LoginDuration.Observe(time.Since(now).Seconds())
	})
}
