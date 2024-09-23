package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ykkalexx/recommendation-system/internal/api"
	"github.com/ykkalexx/recommendation-system/internal/config"
	"github.com/ykkalexx/recommendation-system/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect(client)

	// Check if database is empty and generate fake data if needed
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := client.Database("recommendationDB")
	count, err := db.Collection("behaviors").CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}
	if count == 0 {
		log.Println("Generating fake data...")
		err = database.GenerateFakeData(ctx, db)
		if err != nil {
			log.Fatalf("Failed to generate fake data: %v", err)
		}
		log.Println("Fake data generated successfully")
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