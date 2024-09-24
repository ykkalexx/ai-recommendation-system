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

	// Call Python API for recommendations
	resp, err := http.Get("http://localhost:5000/recommend?user_id=" + userID)
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var recommendations []string
	err = json.NewDecoder(resp.Body).Decode(&recommendations)
	if err != nil {
		http.Error(w, "Failed to process recommendations", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(recommendations)
}