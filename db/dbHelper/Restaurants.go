package dbHelper

import (
	"RMS/db"
	"RMS/models"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
)

func InsertRestaurant(req models.CreateRestaurantRequest, createdBy int64, w http.ResponseWriter) (int64, error) {
	var restaurantID int64

	query := `
		INSERT INTO restaurants 
			(name, address, city, latitude, longitude, rating, is_active, created_by)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := db.RM.QueryRow(
		query,
		req.Name,
		req.Address,
		req.City,
		req.Latitude,
		req.Longitude,
		req.Rating,
		true,
		createdBy,
	).Scan(&restaurantID)

	if err != nil {

		log.Println(err)
		http.Error(w, "failed to insert restaurant: "+err.Error(), http.StatusInternalServerError)

	}

	return restaurantID, nil
}

func GetUserPrimaryCoordinates(userID int64) (float64, float64, error) {
	var lat, lng float64

	query := `
		SELECT latitude, longitude 
		FROM addresses 
		WHERE user_id = $1 AND is_primary = true AND label != 'shop'
		LIMIT 1
	`

	err := db.RM.QueryRow(query, userID).Scan(&lat, &lng)
	if err != nil {
		return 0, 0, errors.New("user's primary address not found")
	}

	return lat, lng, nil
}

// GetRestaurantCoordinates returns restaurant lat, lng, name, and city
func GetRestaurantCoordinates(restaurantID int64) (float64, float64, string, string, error) {
	var lat, lng float64
	var name, city string

	query := `
		SELECT latitude, longitude, name, city 
		FROM restaurants 
		WHERE id = $1
	`

	err := db.RM.QueryRow(query, restaurantID).Scan(&lat, &lng, &name, &city)
	if err != nil {
		return 0, 0, "", "", errors.New("restaurant not found")
	}

	return lat, lng, name, city, nil
}

func AddMenuBySubAdminHelper(req models.AddMenuRequest, restaurantID int64, createdBy int64) (int64, error) {
	var menuID int64

	query := `
		INSERT INTO menus (name, description, price, is_available, food_type, category, restaurant_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := db.RM.QueryRow(query,
		req.Name,
		req.Description,
		req.Price,
		true,
		strings.ToLower(req.FoodType),
		strings.ToLower(req.Category),
		restaurantID,
		createdBy,
	).Scan(&menuID)

	if err != nil {
		return 0, err
	}

	return menuID, nil
}
func CheckRestaurantOwnership(restaurantID int64, userID int64) (bool, error) {
	var count int
	err := db.RM.QueryRow(`
		SELECT COUNT(1) 
		FROM restaurants 
		WHERE id = $1 AND created_by = $2
	`, restaurantID, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func AddMenusByAdminHelper(req []models.AddMenuRequest, restaurantID, userID int64) ([]models.AddedMenu, error) {
	tx, err := db.RM.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	var addedMenus []models.AddedMenu

	validFoodTypes := map[string]bool{"veg": true, "non-veg": true}
	validCategories := map[string]bool{
		"starter": true, "main-course": true, "dessert": true, "breakfast": true, "beverage": true,
	}

	for _, menu := range req {
		if menu.Name == "" || menu.Price < 0 || menu.FoodType == "" || menu.Category == "" {
			err = sql.ErrTxDone
			return nil, err
		}

		if !validFoodTypes[strings.ToLower(menu.FoodType)] {
			err = sql.ErrTxDone
			return nil, err
		}

		if !validCategories[strings.ToLower(menu.Category)] {
			err = sql.ErrTxDone
			return nil, err
		}

		var menuID int64
		query := `
			INSERT INTO menus (name, description, price, food_type, category, restaurant_id, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`
		err = tx.QueryRow(query,
			menu.Name,
			menu.Description,
			menu.Price,
			strings.ToLower(menu.FoodType),
			strings.ToLower(menu.Category),
			restaurantID,
			userID,
		).Scan(&menuID)

		if err != nil {
			return nil, err
		}

		addedMenus = append(addedMenus, models.AddedMenu{
			ID:           menuID,
			Name:         menu.Name,
			Description:  menu.Description,
			Price:        menu.Price,
			IsAvailable:  true,
			FoodType:     strings.ToLower(menu.FoodType),
			Category:     strings.ToLower(menu.Category),
			RestaurantID: restaurantID,
			CreatedBy:    userID,
		})
	}

	return addedMenus, nil
}

func GetAllRestaurantsHelper() ([]models.SlimRestaurantResponse, error) {
	rows, err := db.RM.Query(`
		SELECT id, name, address, city, rating, is_active
		FROM restaurants
		WHERE is_active = true
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []models.SlimRestaurantResponse

	for rows.Next() {
		var rest models.SlimRestaurantResponse
		if err := rows.Scan(
			&rest.ID, &rest.Name, &rest.Address, &rest.City,
			&rest.Rating, &rest.IsActive,
		); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, rest)
	}

	return restaurants, nil
}

// GetRestaurantNameByID checks if restaurant exists & is active, and returns name.
func GetRestaurantNameByID(restaurantID int64) (string, error) {
	var name string
	err := db.RM.QueryRow(`
		SELECT name 
		FROM restaurants 
		WHERE id = $1 AND is_active = true
	`, restaurantID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

// GetMenusByRestaurantID fetches all available menus for given restaurant.
func GetMenusByRestaurantID(restaurantID int64) ([]models.MenuItem, error) {
	rows, err := db.RM.Query(`
		SELECT id, name, description, price, is_available, food_type, category
		FROM menus
		WHERE restaurant_id = $1 AND is_available = true
	`, restaurantID)
	if err != nil {
		return nil, err
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
			continue
		}
		menus = append(menus, item)
	}
	return menus, nil
}
