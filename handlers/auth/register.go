package auth

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req models.RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" || req.Label == "" || req.AddressLine == "" || req.City == "" {
			http.Error(w, "requied filled needed", http.StatusBadRequest)
			return
		}
		if req.Role == "sub-admin" {
			http.Error(w, "sub-admin can only be created by admin , you can create user ", http.StatusBadRequest)
			return
		}
		validLabels := map[string]bool{
			"home":   true,
			"office": true,
			"gym":    true,
			"other":  true,
			"shop":   true,
		}
		if !validLabels[strings.ToLower(req.Label)] {
			http.Error(w, "invalid address label", http.StatusBadRequest)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "failed to start transaction", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
		var exists bool
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email).Scan(&exists)
		if err != nil {
			log.Println("failed to check existing user:", err)
			http.Error(w, "failed to check existing user", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "email already registered", http.StatusConflict)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}
		var userID int64
		err = tx.QueryRow(`
			INSERT INTO users (name, email, password)
			VALUES ($1, $2, $3)
			RETURNING id
		`, req.Name, req.Email, string(hashedPassword)).Scan(&userID)
		if err != nil {
			//log.Println("failed to insert into users:", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
		_, err = tx.Exec(`
			INSERT INTO user_roles (user_id, role)
			VALUES ($1, $2)
		`, userID, req.Role)
		if err != nil {
			//.Println("failed to insert into user_roles:", err)
			http.Error(w, "failed to assign role", http.StatusInternalServerError)
			return
		}
		_, err = tx.Exec(`
			INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
			VALUES ($1, $2, $3, $4, $5, $6, true)
		`, userID, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude)
		if err != nil {
			//log.Println("failed to insert address:", err)
			http.Error(w, "failed to insert address", http.StatusInternalServerError)
			return
		}
		if err = tx.Commit(); err != nil {
			//log.Println("failed to commit transaction:", err)
			http.Error(w, "failed to complete registration", http.StatusInternalServerError)
			return
		}
		accessToken, err := utils.GenerateAccessToken(userID, req.Email, req.Role)
		if err != nil {
			log.Println("failed to generate access token:", err)
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(userID, req.Email, req.Role)
		if err != nil {
			log.Println("failed to generate refresh token:", err)
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}
		resp := models.RegisterResponse{
			Message: "user registered successfully",
			User: models.RegisteredUser{
				ID:    userID,
				Name:  req.Name,
				Email: req.Email,
				Role:  req.Role,
			},
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)

	}
}
