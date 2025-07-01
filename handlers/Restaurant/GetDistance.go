package restaurant

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
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
			http.Error(w, "invalid restaurant ID", http.StatusBadRequest)
			return
		}
		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		userID := claims.UserID
		role := claims.Role
		if role != "user" {
			http.Error(w, "unauthorized : soory if you are user tehn only you should use this api", http.StatusUnauthorized)
		}

		var userLat, userLng float64
		err = db.QueryRow(`
			SELECT latitude, longitude 
			FROM addresses 
			WHERE user_id = $1 AND is_primary = true AND label != 'shop'
			LIMIT 1
		`, userID).Scan(&userLat, &userLng)
		if err != nil {
			http.Error(w, "user's primary address not found", http.StatusNotFound)
			return
		}
		var restLat, restLng float64
		var restName, restCity string
		err = db.QueryRow(`
			SELECT latitude, longitude, name, city 
			FROM restaurants 
			WHERE id = $1
		`, restaurantID).Scan(&restLat, &restLng, &restName, &restCity)
		if err != nil {
			http.Error(w, "restaurant not found", http.StatusNotFound)
			return
		}
		distanceKM := utils.CalculateDistance(userLat, userLng, restLat, restLng)
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
