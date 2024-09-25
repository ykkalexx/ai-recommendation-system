package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/ykkalexx/recommendation-system/internal/database"
	"github.com/ykkalexx/recommendation-system/internal/models"
	"github.com/ykkalexx/recommendation-system/internal/utils"
)

type RecommendationsResponse struct {
	Recommendations []string `json:"recommendations"`
	UserId string `json:"user_id"`
}


func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/api/behavior", recordBehavior).Methods("POST")
	r.HandleFunc("/api/recommendations/{user_id}", getRecommendations).Methods("GET")
}

func recordBehavior(w http.ResponseWriter, r *http.Request) {
	var behavior models.UserBehavior
	err := json.NewDecoder(r.Body).Decode(&behavior)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to decode behavior: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("recommendationDB", "behaviors")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, behavior)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to record behavior: %v", err)
		http.Error(w, "Failed to record behavior", http.StatusInternalServerError)
		return
	}

	utils.InfoLogger.Printf("Behavior recorded for user %s", behavior.UserID)
	w.WriteHeader(http.StatusCreated)
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]

    // Call Python API for recommendations
    resp, err := http.Get("http://127.0.0.1:5000/recommend?user_id=" + userID)
    if err != nil {
		utils.ErrorLogger.Printf("Failed to get recommendations from ML API: %v", err)
        http.Error(w, "Failed to get recommendations: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        http.Error(w, "Python API returned status code: "+strconv.Itoa(resp.StatusCode), http.StatusInternalServerError)
        return
    }

	var response RecommendationsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to process recommendations: %v", err)
		http.Error(w, "Failed to process recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	utils.InfoLogger.Printf("Recommendations retrieved for user %s", userID)
	json.NewEncoder(w).Encode(response.Recommendations)
}