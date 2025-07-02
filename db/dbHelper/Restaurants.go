package dbHelper

import (
	"RMS/db"
	"RMS/models"
	"database/sql"
	"net/http"
)

func CheckRestaurant(restaurantID int64, w http.ResponseWriter) error {
	var restaurantExists bool
	err := db.RM.QueryRow(`SELECT EXISTS(SELECT 1 FROM restaurants WHERE id = $1)`, restaurantID).Scan(&restaurantExists)
	if err != nil || !restaurantExists {
		http.Error(w, "restaurant does not exist", http.StatusBadRequest)
		return err
	}
	return nil
}

// fucn to check the user exist or not

func CheckUserExsist(UserId int64, w http.ResponseWriter) error {
	var userExists bool
	err := db.RM.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, UserId).Scan(&userExists)
	if err != nil || !userExists {
		http.Error(w, "invalid user", http.StatusBadRequest)
		return err
	}
	return nil
}

func AddMenuInDb(tx *sql.Tx, req *models.AddNewAddressRequest, w http.ResponseWriter) error {

	return nil

}

func GetUserLatLong(userID int64, w http.ResponseWriter) error  {
	var userLat, userLng float64
	err := db.RM.QueryRow(`
			SELECT latitude, longitude 
			FROM addresses 
			WHERE user_id = $1 AND is_primary = true AND label != 'shop'
			LIMIT 1
		`, userID).Scan(&userLat, &userLng)
	if err != nil {
		http.Error(w, "user's primary address not found", http.StatusNotFound)
		return err
	}
	return nil
}

func GetRestrLatLong(restaurantID int64 ,w http.ResponseWriter) error lat , long  {
	var restLat, restLng float64
	err := db.RM.QueryRow(`
			SELECT latitude, longitude
			FROM restaurants 
			WHERE id = $1
		`, restaurantID).Scan(&restLat, &restLng)
	if err != nil {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return err
	}
	return nil
}
func CreateRestaurantBySubAdmin(restaurantID int64, w http.ResponseWriter) RestaurantId int64 , error {
var restaurantID int64
query := `
			INSERT INTO restaurants (name, address, city, latitude, longitude, rating, is_active, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
         err = db.RM.QueryRow(query,
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
http.Error(w, "failed to create restaurant", http.StatusInternalServerError)
log.Println("insert error:", err)
return
}


}