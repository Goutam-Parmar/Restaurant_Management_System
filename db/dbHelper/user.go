package dbHelper

import (
	"RMS/db"
	"RMS/models"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func SearchEmail(tx *sql.Tx, req *models.RegisterRequestDB, w http.ResponseWriter) error {
	var emailExists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email).Scan(&emailExists)
	if err != nil {
		return err
	}
	if emailExists {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("Email already exists")

	}
	return err
}

func SearchPhone(tx *sql.Tx, req *models.RegisterRequestDB, w http.ResponseWriter) error {
	var phoneExists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`, req.Phone).Scan(&phoneExists)
	if err != nil {
		return err
	}
	if phoneExists {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("Phone already exists")

	}
	return err

}

// list of all the users
func GetAllSubAdmin(w http.ResponseWriter) ([]models.UserBrief2, error) {

	rows, err := db.RM.Query(`
			         SELECT 
                     u.id, 
                     u.name, 
                     a.city
             FROM users u
             JOIN user_roles ur ON u.id = ur.user_id
             JOIN addresses a ON u.id = a.user_id AND a.is_primary = true
        WHERE ur.role = 'sub-admin';
		`)
	if err != nil {
		http.Error(w, "Failed to fetch subadmins", http.StatusInternalServerError)
		return nil, err
	}

	defer rows.Close()
	var users []models.UserBrief2
	for rows.Next() {
		var u models.UserBrief2
		if err := rows.Scan(&u.ID, &u.Name, &u.City); err != nil {

			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// lsit all user
func GetListOfUsers() ([]models.UserBrief2, error) {
	query := `
		SELECT u.id, u.name, a.city
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN addresses a ON u.id = a.user_id AND a.is_primary = true
		WHERE ur.role = 'user'
	`
	rows, err := db.RM.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserBrief2
	for rows.Next() {
		var user models.UserBrief2
		if err := rows.Scan(&user.ID, &user.Name, &user.City); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// dbHelper/user.go

func RegisterUser(tx *sql.Tx, req *models.RegisterRequestDB) error {
	// Insert into users
	query := `
		INSERT INTO users (name, email, password, phone, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := tx.QueryRow(query, req.Name, req.Email, req.Password, req.Phone, req.CreatedBy).Scan(&req.UserId)
	if err != nil {
		return err
	}

	// Insert into user_roles
	_, err = tx.Exec(`
		INSERT INTO user_roles (user_id, role)
		VALUES ($1, $2)
	`, req.UserId, req.Role)
	if err != nil {
		return err
	}

	// Insert into addresses
	_, err = tx.Exec(`
		INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		req.UserId,
		req.Label,
		req.AddressLine,
		req.City,
		req.Latitude,
		req.Longitude,
		true,
	)
	return err
}

// to check login with email
func CheckLoginCredentials(req *models.LoginRequest, w http.ResponseWriter) error {
	var userID int64
	var name, hashedPassword, role string
	err := db.RM.QueryRow(`
			SELECT u.id, u.name, u.password, ur.role
			FROM users u
			JOIN user_roles ur ON u.id = ur.user_id
			WHERE u.email = $1
			LIMIT 1
		`, req.Email).Scan(&userID, &name, &hashedPassword, &role)
	req.UserId = userID
	req.Role = role
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return err
	}
	return err

}

// to add new add to user
func AddNewAddress(tx *sql.Tx, req *models.AddNewAddressRequest, w http.ResponseWriter) error {
	//var id int64
	if req.IsPrimary {
		_, err := tx.Exec(`UPDATE addresses SET is_primary = false WHERE user_id = $1`, req.UserId)
		if err != nil {
			http.Error(w, "failed to unset other primary addresses", http.StatusInternalServerError)
			return err
		}
	}

	var exists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM addresses WHERE user_id = $1 AND label = $2)`, req.UserId, req.Label).Scan(&exists)
	if err != nil {
		http.Error(w, "failed to check existing address", http.StatusInternalServerError)
		return err
	}
	if exists {
		_, err = tx.Exec(`
				UPDATE addresses
				SET address_line = $1, city = $2, latitude = $3, longitude = $4, is_primary = $5
				WHERE user_id = $6 AND label = $7
			`,
			req.AddressLine, req.City, req.Latitude, req.Longitude, req.IsPrimary, req.UserId, req.Label,
		)
	} else {
		// Insert
		_, err = tx.Exec(`
				INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`,
			req.UserId, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude, req.IsPrimary,
		)
	}
	if err != nil {
		http.Error(w, "failed to upsert address: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
