package restaurant

import (
	"RMS/models"
	"RMS/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func CreateRestaurant(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req models.CreateRestaurantRequest

		body, _ := io.ReadAll(r.Body)
		log.Println("Raw Request Body:", string(body))
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
			claims.UserID,
		).Scan(&restaurantID)

		if err != nil {
			http.Error(w, "Failed to create restaurant", http.StatusInternalServerError)
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
				CreatedBy: claims.UserID,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
