package auth

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func CreateSubAdmin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized: missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		_, claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}
		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "forbidden: only admins can create sub-admins", http.StatusForbidden)
			return
		}
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		if req.Name == "" || req.Email == "" || req.Password == "" || req.Label == "" || req.AddressLine == "" || req.City == "" {
			http.Error(w, " required fields", http.StatusBadRequest)
			return
		}
		validLabels := map[string]bool{"home": true, "office": true, "gym": true, "other": true, "shop": true}
		if !validLabels[strings.ToLower(req.Label)] {
			http.Error(w, "invalid address label", http.StatusBadRequest)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "DB transaction error", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback()
			}
		}()

		var exists bool
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email).Scan(&exists)
		if err != nil {
			http.Error(w, " DB check error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Password hash failed", http.StatusInternalServerError)
			return
		}
		var userID int64
		err = tx.QueryRow(`
			INSERT INTO users (name, email, password) 
			VALUES ($1, $2, $3) RETURNING id
		`, req.Name, req.Email, string(hashedPassword)).Scan(&userID)
		if err != nil {
			http.Error(w, "insert user failed", http.StatusInternalServerError)
			return
		}
		_, err = tx.Exec(`
			INSERT INTO user_roles (user_id, role) 
			VALUES ($1, 'sub-admin')
		`, userID)
		if err != nil {
			http.Error(w, "insert role failed", http.StatusInternalServerError)
			return
		}
		_, err = tx.Exec(`
			INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
			VALUES ($1, $2, $3, $4, $5, $6, true)
		`, userID, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude)
		if err != nil {
			http.Error(w, "insert address failed", http.StatusInternalServerError)
			return
		}
		if err = tx.Commit(); err != nil {
			http.Error(w, "commit failed", http.StatusInternalServerError)
			return
		}
		resp := models.CreateUserResponse{
			Message:        "Sub-admin created successfully",
			UserID:         userID,
			Email:          req.Email,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	}
}
