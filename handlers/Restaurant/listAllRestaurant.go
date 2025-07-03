package restaurant

import (
	"RMS/db/dbHelper"
	"RMS/models"
	"encoding/json"
	"net/http"
	"time"
)

func GetAllRestaurants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		restaurants, err := dbHelper.GetAllRestaurantsHelper()
		if err != nil {
			http.Error(w, "failed to fetch restaurants", http.StatusInternalServerError)
			return
		}

		resp := models.AllRestaurantsSlimResponse{
			Message:        "Restaurants fetched",
			Restaurants:    restaurants,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
