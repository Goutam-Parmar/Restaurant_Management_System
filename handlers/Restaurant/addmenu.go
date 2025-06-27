package restaurant

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddMenu(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		vars := mux.Vars(r)
		restaurantIDStr := vars["restaurant_id"]
		log.Println("extracted restaurant_id sucess:", restaurantIDStr)
		restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)
		if err != nil {
			//log.Println("invalid restaurant_id:", restaurantIDStr, "| Error:", err)
			http.Error(w, "invalid restaurant_id", http.StatusBadRequest)
			return
		}
		var req models.AddMenuRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("failed to decode request body:", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Name == "" || req.Price < 0 || req.FoodType == "" || req.Category == "" {
			http.Error(w, "missing required fields", http.StatusBadRequest)
			return
		}
		validFoodTypes := map[string]bool{"veg": true, "non-veg": true}
		validCategories := map[string]bool{
			"starter": true, "main-course": true, "dessert": true, "breakfast": true, "beverage": true,
		}
		if !validFoodTypes[strings.ToLower(req.FoodType)] {
			http.Error(w, "invalid food_type", http.StatusBadRequest)
			return
		}
		if !validCategories[strings.ToLower(req.Category)] {
			http.Error(w, " invalid category", http.StatusBadRequest)
			return
		}
		tokenStr := r.Header.Get("Authorization")
		userID, email, err := utils.ExtractUserIDAndEmailFromHeader(tokenStr)
		if err != nil {
			log.Println("failed to extract user from token:", err)
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		var menuID int64
		err = db.QueryRow(`
			INSERT INTO menus (name, description, price, food_type, category, restaurant_id, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`, req.Name, req.Description, req.Price, req.FoodType, req.Category, restaurantID, userID).Scan(&menuID)

		if err != nil {
			http.Error(w, "failed to add menu item", http.StatusInternalServerError)
			return
		}
		resp := models.AddMenuResponse{
			Message: "Menu item added successfully",
			Menu: models.AddedMenu{
				ID:           menuID,
				Name:         req.Name,
				Description:  req.Description,
				Price:        req.Price,
				IsAvailable:  true,
				FoodType:     req.FoodType,
				Category:     req.Category,
				RestaurantID: restaurantID,
				CreatedBy:    userID,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
