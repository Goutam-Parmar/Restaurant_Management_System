package restaurant

import (
	"RMS/db/dbHelper"
	"RMS/models"
	"encoding/json"
	"github.com/gorilla/mux"
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

		restaurantName, err := dbHelper.GetRestaurantNameByID(restaurantID)
		if err != nil {
			http.Error(w, "restaurant not found or inactive", http.StatusNotFound)
			return
		}

		menus, err := dbHelper.GetMenusByRestaurantID(restaurantID)
		if err != nil {
			http.Error(w, "failed to fetch menu items", http.StatusInternalServerError)
			return
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
