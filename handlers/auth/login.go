package auth

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" || req.Password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}
		var userID int64
		var name, hashedPassword, role string
		err := db.QueryRow(`
			SELECT u.id, u.name, u.password, ur.role
			FROM users u
			JOIN user_roles ur ON u.id = ur.user_id
			WHERE u.email = $1
			LIMIT 1
		`, req.Email).Scan(&userID, &name, &hashedPassword, &role)
		if err != nil {
			log.Println(err)
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		accessToken, err := utils.GenerateAccessToken(userID, req.Email, role)
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(userID, req.Email, role)
		if err != nil {
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}
		resp := models.LoginResponse{
			Message:      "Login successful",
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User: models.LoginUser{
				ID:    userID,
				Name:  name,
				Email: req.Email,
				Role:  role,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	}
}
