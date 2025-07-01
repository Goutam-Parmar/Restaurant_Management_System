package restaurant

import (
	"RMS/db"
	"RMS/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetMenusByRestaurantID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		vars := mux.Vars(r)
		restaurantIDStr := vars["restaurant_id"]
		restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid restaurant ID", http.StatusBadRequest)
			return
		}
		var restaurantName string
		err = db.RM.QueryRow("SELECT name FROM restaurants WHERE id = $1 AND is_active = true", restaurantID).Scan(&restaurantName)
		if err != nil {
			http.Error(w, "restaurant not found or inactive", http.StatusNotFound)
			return
		}
		rows, err := db.RM.Query(`
			SELECT id, name, description, price, is_available, food_type, category
			FROM menus
			WHERE restaurant_id = $1 AND is_available = true
		`, restaurantID)
		if err != nil {
			http.Error(w, "failed to fetch menu items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var menus []models.MenuItem
		for rows.Next() {
			var item models.MenuItem
			if err := rows.Scan(
				&item.ID,
				&item.Name,
				&item.Description,
				&item.Price,
				&item.IsAvailable,
				&item.FoodType,
				&item.Category,
			); err != nil {
				log.Println("Row scan error:", err)
				continue
			}
			menus = append(menus, item)
		}
		resp := models.MenuListResponse{
			Message:        "Menus fetched successfully",
			RestaurantID:   restaurantID,
			RestaurantName: restaurantName,
			Menus:          menus,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
