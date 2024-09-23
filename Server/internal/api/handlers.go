package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ykkalexx/recommendation-system/internal/database"
	"github.com/ykkalexx/recommendation-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/api/behavior", recordBehavior).Methods("POST")
	r.HandleFunc("/api/recommendations/{user_id}", getRecommendations).Methods("GET")
}

func recordBehavior(w http.ResponseWriter, r *http.Request) {
	var behavior models.UserBehavior
	err := json.NewDecoder(r.Body).Decode(&behavior)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("recommendationDB", "behaviors")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, behavior)
	if err != nil {
		log.Printf("Error inserting behavior: %v", err)
		http.Error(w, "Failed to record behavior", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	collection := database.GetCollection("recommendationDB", "behaviors")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For now, we'll just return the last 5 items the user interacted with
	// In a real system, this would be replaced with actual recommendation logic
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(5))
	if err != nil {
		log.Printf("Error fetching recommendations: %v", err)
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var behaviors []models.UserBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		log.Printf("Error decoding recommendations: %v", err)
		http.Error(w, "Failed to process recommendations", http.StatusInternalServerError)
		return
	}

	recommendations := make([]string, len(behaviors))
	for i, behavior := range behaviors {
		recommendations[i] = behavior.ItemID
	}

	json.NewEncoder(w).Encode(recommendations)
}