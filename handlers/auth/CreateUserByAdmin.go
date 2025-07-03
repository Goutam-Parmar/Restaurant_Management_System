package auth

import (
	"RMS/db"
	"RMS/db/dbHelper"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser() http.HandlerFunc {
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
		var req models.RegisterRequestDB
		req.CreatedBy = authClaims.UserID

		if authClaims.Role != "admin" {
			http.Error(w, "this is api for admin", http.StatusUnauthorized)
			return
		}
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
		//pasword
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if hashErr != nil {
			http.Error(w, "password hash error", http.StatusBadRequest)
			return
		}
		req.Password = string(hashedPassword)
		//since it is a mutiple table insert we will use transaction
		tx, err := db.RM.Begin()
		defer utils.Tx(tx, &err)
		err = dbHelper.SearchEmail(tx, &req, w)
		err = dbHelper.SearchPhone(tx, &req, w)
		err = dbHelper.RegisterUser(tx, &req)
		if err = tx.Commit(); err != nil {
			http.Error(w, "failed to complete registration", http.StatusInternalServerError)
			return
		}
		accessToken, err := utils.GenerateAccessToken(req.UserId, req.Email, req.Role)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(req.UserId, req.Email, req.Role)
		if err != nil {
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}
		resp := models.RegisterResponse{
			User: models.RegisteredUser{
				ID:    req.UserId,
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
