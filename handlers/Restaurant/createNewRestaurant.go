package restaurant

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func CreateRestaurant(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req models.CreateRestaurantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			log.Println("failed to decode request:", err)
			return
		}
		authHeader := r.Header.Get("Authorization")
		userID, _, err := utils.ExtractUserIDAndEmailFromHeader(authHeader)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			log.Println("failed to extract user from token:", err)
			return
		}
		var restaurantID int64
		query := `
			INSERT INTO restaurants (name, address, city, latitude, longitude, rating, is_active, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err = db.QueryRow(query,
			req.Name,
			req.Address,
			req.City,
			req.Latitude,
			req.Longitude,
			req.Rating,
			true,
			userID,
		).Scan(&restaurantID)
		if err != nil {
			http.Error(w, "failed to create restaurant", http.StatusInternalServerError)
			log.Println("db insert error:", err)
			return
		}
		resp := models.CreateRestaurantResponse{
			Message: "Restaurant created successfully",
			Restaurant: models.CreatedRestaurant{
				ID:        restaurantID,
				Name:      req.Name,
				Address:   req.Address,
				City:      req.City,
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
				Rating:    req.Rating,
				IsActive:  true,
				CreatedBy: userID,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)

	}
}
