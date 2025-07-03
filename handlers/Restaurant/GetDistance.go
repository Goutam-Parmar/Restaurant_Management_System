package restaurant

import (
	"RMS/db/dbHelper"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func GetDistanceToRestaurant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		vars := mux.Vars(r)
		restaurantIDStr := vars["restaurant_id"]
		restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
			return
		}
		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil || claims.Role != "user" {
			http.Error(w, "Unauthorized: only users can use this API", http.StatusUnauthorized)
			return
		}
		userLat, userLong, err := dbHelper.GetUserPrimaryCoordinates(claims.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		restLat, restLong, restName, restCity, err := dbHelper.GetRestaurantCoordinates(restaurantID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		distanceKM := utils.CalculateDistance(userLat, userLong, restLat, restLong)
		resp := models.RestaurantDistanceResponse{
			Message:    "Distance calculated successfully",
			DistanceKM: distanceKM,
			Restaurant: models.MinimalRestaurant{
				ID:   restaurantID,
				Name: restName,
				City: restCity,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
