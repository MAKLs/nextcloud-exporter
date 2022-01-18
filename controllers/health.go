package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthController serves the status of the exporter.
func HealthController() http.Handler {
	// TODO: provide some stats about the exporter (i.e. last scrape time)
	health := struct {
		Status string `json:"status"`
	}{Status: "UP"}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(health)
		if err != nil {
			log.Fatalf("failed to serialize health status: %s", err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(body)
	})
}
