package restaurant

import (
	"RMS/db/dbHelper"
	"RMS/models"
	"RMS/utils"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func CreateRestaurantByAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req models.CreateRestaurantRequest

		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(" Failed to decode CreateRestaurantRequest:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(w, "Forbidden: only admins can create restaurants", http.StatusForbidden)
			return
		}
		if req.Rating < 0.0 || req.Rating > 5.0 {
			http.Error(w, "Rating must be between 0.0 and 5.0", http.StatusBadRequest)
			return
		}

		restaurantID, err := dbHelper.InsertRestaurant(req, claims.UserID, w)
		resp := models.CreateRestaurantResponse{
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
