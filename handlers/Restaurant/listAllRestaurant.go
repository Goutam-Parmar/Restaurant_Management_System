package restaurant

import (
	"RMS/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func GetAllRestaurants(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rows, err := db.Query(`
			SELECT id, name, address, city, rating, is_active
			FROM restaurants
			WHERE is_active = true
			ORDER BY name
		`)
		if err != nil {
			http.Error(w, "failed to fetch restaurants", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var restaurants []models.SlimRestaurantResponse
		for rows.Next() {
			var rest models.SlimRestaurantResponse
			if err := rows.Scan(
				&rest.ID, &rest.Name, &rest.Address, &rest.City,
				&rest.Rating, &rest.IsActive,
			); err != nil {
				http.Error(w, "Error reading data", http.StatusInternalServerError)
				return
			}
			restaurants = append(restaurants, rest)
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
