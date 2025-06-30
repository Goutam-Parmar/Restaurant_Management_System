package auth

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}

		authClaims, err := utils.ExtractAuthClaims(authHeader)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		createdByID := authClaims.UserID

		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
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
			http.Error(w, "required fields missing", http.StatusBadRequest)
			return
		}

		if len(req.Phone) != 10 || !utils.IsNumeric(req.Phone) {
			http.Error(w, "invalid phone number", http.StatusBadRequest)
			return
		}

		if req.Role == "sub-admin" && authClaims.Role != "admin" {
			http.Error(w, "only admin can create sub-admins", http.StatusForbidden)
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

		var emailExists bool
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email).Scan(&emailExists)
		if err != nil {
			http.Error(w, "failed to check existing email", http.StatusInternalServerError)
			return
		}
		if emailExists {
			http.Error(w, "email already registered", http.StatusConflict)
			return
		}

		// Check phone exists
		var phoneExists bool
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`, req.Phone).Scan(&phoneExists)
		if err != nil {
			http.Error(w, "failed to check existing phone", http.StatusInternalServerError)
			return
		}
		if phoneExists {
			http.Error(w, "phone already registered", http.StatusConflict)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		var userID int64
		err = tx.QueryRow(`
			INSERT INTO users (name, email, password, phone, created_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, req.Name, req.Email, string(hashedPassword), req.Phone, createdByID).Scan(&userID)
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`INSERT INTO user_roles (user_id, role) VALUES ($1, $2)`, userID, req.Role)
		if err != nil {
			http.Error(w, "failed to assign role", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`
			INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
			VALUES ($1, $2, $3, $4, $5, $6, TRUE)
		`, userID, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude)
		if err != nil {
			http.Error(w, "failed to insert address", http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			http.Error(w, "failed to complete registration", http.StatusInternalServerError)
			return
		}

		accessToken, err := utils.GenerateAccessToken(userID, req.Email, req.Role)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := utils.GenerateRefreshToken(userID, req.Email, req.Role)
		if err != nil {
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
