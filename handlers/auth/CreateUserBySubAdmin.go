package auth

import (
	"RMS/db"
	"RMS/models"
	"RMS/utils"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func CreateUserBySubAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}
		if claims.Role != "sub-admin" {
			http.Error(w, "unauthorized: only sub-admins can create users", http.StatusUnauthorized)
			return
		}
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.Password = strings.TrimSpace(req.Password)
		req.Phone = strings.TrimSpace(req.Phone)
		req.Role = strings.TrimSpace(req.Role)
		req.Label = strings.TrimSpace(req.Label)
		req.AddressLine = strings.TrimSpace(req.AddressLine)
		req.City = strings.TrimSpace(req.City)

		if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" ||
			req.Role == "" || req.Label == "" || req.AddressLine == "" || req.City == "" {
			http.Error(w, "required fields are missing", http.StatusBadRequest)
			return
		}
		if len(req.Phone) != 10 {
			http.Error(w, "phone number must be 10 digits", http.StatusBadRequest)
			return
		}
		if req.Role == "admin" || req.Role == "sub-admin" {
			http.Error(w, "unauthorized: sub-admins can only create users", http.StatusUnauthorized)
			return
		}

		validLabels := map[string]bool{"home": true, "office": true, "gym": true, "other": true, "shop": true}
		if !validLabels[strings.ToLower(req.Label)] {
			http.Error(w, "invalid address label", http.StatusBadRequest)
			return
		}

		// Begin transaction
		tx, err := db.RM.Begin()
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
			http.Error(w, "DB check error", http.StatusInternalServerError)
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
			INSERT INTO users (name, email, password, phone, created_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, req.Name, req.Email, string(hashedPassword), req.Phone, claims.UserID).Scan(&userID)
		if err != nil {
			http.Error(w, "failed to insert user", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`
			INSERT INTO user_roles (user_id, role)
			VALUES ($1, $2)
		`, userID, req.Role)
		if err != nil {
			http.Error(w, "failed to assign role", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`
			INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
			VALUES ($1, $2, $3, $4, $5, $6, true)
		`, userID, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude)
		if err != nil {
			http.Error(w, "failed to insert address", http.StatusInternalServerError)
			return
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			http.Error(w, "commit failed", http.StatusInternalServerError)
			return
		}

		resp := models.CreateUserResponse{
			Message:        "User created successfully by sub-admin",
			UserID:         userID,
			Email:          req.Email,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
