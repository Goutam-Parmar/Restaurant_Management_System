package restaurant

import (
	"RMS/db"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddMenuBySubadmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//Extract restaurant_id
		params := mux.Vars(r)
		restaurantIDStr := params["restaurant_id"]
		restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
			return
		}

		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if claims.Role != "sub-admin" {
			http.Error(w, "Unauthorized: only sub-admins can add menus", http.StatusForbidden)
			return
		}

		var count int
		err = db.RM.QueryRow(`
			SELECT COUNT(1) 
			FROM restaurants 
			WHERE id = $1 AND created_by = $2
		`, restaurantID, claims.UserID).Scan(&count)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			http.Error(w, "You are not allowed to modify this restaurant, you can only add menu which restorant belong to you", http.StatusForbidden)
			return
		}

		var req models.AddMenuRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var menuID int64
		query := `
			INSERT INTO menus (name, description, price, is_available, food_type, category, restaurant_id, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err = db.RM.QueryRow(query,
			req.Name,
			req.Description,
			req.Price,
			true,
			strings.ToLower(req.FoodType),
			strings.ToLower(req.Category),
			restaurantID,
			claims.UserID,
		).Scan(&menuID)
		if err != nil {
			http.Error(w, "Failed to add menu", http.StatusInternalServerError)
			return
		}

		resp := models.AddMenuResponseSubAdmin{
			Message:        "Menu added successfully",
			MenuID:         menuID,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
