package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ykkalexx/recommendation-system/internal/api"
	"github.com/ykkalexx/recommendation-system/internal/models"
)

func TestRecordBehavior(t *testing.T) {
    behavior := models.UserBehavior{
        UserID: "test_user",
        // fill in other fields
    }
    behaviorBytes, _ := json.Marshal(behavior)

    req, err := http.NewRequest("POST", "/api/behavior", bytes.NewBuffer(behaviorBytes))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(api.RecordBehavior)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
    }
}

func TestGetRecommendations(t *testing.T) {
    req, err := http.NewRequest("GET", "/api/recommendations/test_user", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    router := mux.NewRouter()
    router.HandleFunc("/api/recommendations/{user_id}", api.GetRecommendations)

    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}