package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/ykkalexx/recommendation-system/internal/models"
)

var logger = logrus.New()

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

	// TODO: Save behavior to database
	// For now, we'll just log it
	logger.Info("Recorded behavior: %+v", behavior)

	w.WriteHeader(http.StatusCreated)
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		http.Error(w, "user_id is missing in parameters", http.StatusBadRequest)
        return
	}
	
	// TODO: Get recommendations from ML model
	// For now, we'll return mock data
	recommendations := []string{"item1", "item2", "item3"}

	logger.Info("Fetching recommendations for user: %s", userID)

	json.NewEncoder(w).Encode(recommendations)
}