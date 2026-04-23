package main

import (
	"log"
	"net/http"

	"github.com/jemiezler/Vers/backend/internal/config"
	dashboardhttp "github.com/jemiezler/Vers/backend/internal/dashboard/http"
	"github.com/jemiezler/Vers/backend/internal/review"
)

func main() {
	cfg := config.Load()
	service, err := review.NewServiceFromConfig(cfg)
	if err != nil {
		log.Fatalf("configure review service: %v", err)
	}
	handler := dashboardhttp.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.Health)
	mux.HandleFunc("POST /reviews", handler.CreateReview)

	log.Printf("vers api listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
		log.Fatal(err)
	}
}
