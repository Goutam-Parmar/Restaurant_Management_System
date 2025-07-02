package restaurant

import (
	"RMS/db"
	"RMS/db/dbHelper"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func CreateRestaurantBySubAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var req models.CreateRestaurantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			log.Println("failed to decode request:", err)
			return
		}

		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if strings.ToLower(claims.Role) != "sub-admin" {
			http.Error(w, "forbidden: Only sub-admins can create restaurants here", http.StatusForbidden)
			return
		}
		err := dbHelper.CreateRestaurantBySubAdmin()

		resp := models.CreateRestaurantResponse{
			Message: "restaurant created by sub-admin successfully",
			Restaurant: models.CreatedRestaurant{
				ID:        restaurantID,
				Name:      req.Name,
				Address:   req.Address,
				City:      req.City,
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
				Rating:    req.Rating,
				IsActive:  true,
				CreatedBy: claims.UserID,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
