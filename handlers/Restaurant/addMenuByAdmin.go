package restaurant

import (
	"RMS/db"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddMenuByAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		vars := mux.Vars(r)
		restaurantIDStr := vars["restaurant_id"]
		restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid restaurant_id", http.StatusBadRequest)
			return
		}

		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			log.Println("failed to extract auth claims:", err)
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		userID := claims.UserID
		log.Println("Authenticated User ID:", userID)

		var req []models.AddMenuRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("failed to decode request body:", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		var restaurantExists bool
		err = db.RM.QueryRow(`SELECT EXISTS(SELECT 1 FROM restaurants WHERE id = $1)`, restaurantID).Scan(&restaurantExists)
		if err != nil || !restaurantExists {
			http.Error(w, "restaurant does not exist", http.StatusBadRequest)
			return
		}

		// Validate user exists
		var userExists bool
		err = db.RM.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&userExists)
		if err != nil || !userExists {
			http.Error(w, "invalid user", http.StatusBadRequest)
			return
		}

		// Prepare for insertion
		validFoodTypes := map[string]bool{"veg": true, "non-veg": true}
		validCategories := map[string]bool{
			"starter": true, "main-course": true, "dessert": true, "breakfast": true, "beverage": true,
		}

		tx, err := db.RM.Begin()
		if err != nil {
			log.Println("failed to begin tx:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		var addedMenus []models.AddedMenu

		for _, menu := range req {
			// Basic validation
			if menu.Name == "" || menu.Price < 0 || menu.FoodType == "" || menu.Category == "" {
				err = tx.Rollback()
				http.Error(w, "missing required fields in menu item: "+menu.Name, http.StatusBadRequest)
				return
			}

			if !validFoodTypes[strings.ToLower(menu.FoodType)] {
				err = tx.Rollback()
				http.Error(w, "invalid food_type for "+menu.Name, http.StatusBadRequest)
				return
			}

			if !validCategories[strings.ToLower(menu.Category)] {
				err = tx.Rollback()
				http.Error(w, "invalid category for "+menu.Name, http.StatusBadRequest)
				return
			}

			var menuID int64
			err = tx.QueryRow(`
				INSERT INTO menus (name, description, price, food_type, category, restaurant_id, created_by)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
			`,
				menu.Name,
				menu.Description,
				menu.Price,
				menu.FoodType,
				menu.Category,
				restaurantID,
				userID,
			).Scan(&menuID)

			if err != nil {
				log.Printf("failed to insert menu: %+v | err: %v", menu, err)
				http.Error(w, "failed to add menu item: "+err.Error(), http.StatusInternalServerError)
				return
			}

			addedMenus = append(addedMenus, models.AddedMenu{
				ID:           menuID,
				Name:         menu.Name,
				Description:  menu.Description,
				Price:        menu.Price,
				IsAvailable:  true,
				FoodType:     menu.FoodType,
				Category:     menu.Category,
				RestaurantID: restaurantID,
				CreatedBy:    userID,
			})
		}

		resp := models.AddMenuResponse{
			Message:        "Menus added successfully",
			Menus:          addedMenus,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
