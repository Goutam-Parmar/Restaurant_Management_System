package dbHelper

import (
	"RMS/models"
	"database/sql"
)

// Return the created user ID
func RegisterUser(tx *sql.Tx, req *models.RegisterRequestDB) error {
	var userId int64
	err := tx.QueryRow(`INSERT INTO users (name, email, password, phone, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`, req.Name, req.Email, req.Password, req.Phone, req.Created_by).Scan(&userId)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO user_roles (user_id, role) VALUES ($1, $2)`, userId, req.Role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
			VALUES ($1, $2, $3, $4, $5, $6, TRUE)
		`, userId, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude)
	if err != nil {
		return err
	}
	return nil
}
