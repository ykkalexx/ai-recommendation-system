package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ykkalexx/recommendation-system/internal/api"
	"github.com/ykkalexx/recommendation-system/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// create router
	r := mux.NewRouter()

	// setting up routes
	api.SetupRoutes(r)

	// start server
	addr := cfg.ServerAddress + ":" + cfg.ServerPort
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}