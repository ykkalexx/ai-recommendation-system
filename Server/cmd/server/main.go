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
	"github.com/ykkalexx/recommendation-system/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to connect to Mongodb: %v", err)
	}
	defer database.Disconnect(client)

	// Check if database is empty and generate fake data if needed
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := client.Database("recommendationDB")
	count, err := db.Collection("behaviors").CountDocuments(ctx, bson.M{})
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to count documents: %v", err)
	}
	if count == 0 {
		log.Println("Generating fake data...")
		err = database.GenerateFakeData(ctx, db)
		if err != nil {
			utils.ErrorLogger.Fatalf("Failed to generate data: %v", err)
		}
		utils.InfoLogger.Println("Fake data generated successfully")
	}

	// create router
	r := mux.NewRouter()

	// setting up routes
	api.SetupRoutes(r)

	// start server
	addr := cfg.ServerAddress + ":" + cfg.ServerPort
	utils.InfoLogger.Printf("Server starting on %s...", addr)
	utils.ErrorLogger.Fatal(http.ListenAndServe(addr, r))
}