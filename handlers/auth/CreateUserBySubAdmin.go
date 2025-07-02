package auth

import (
	"RMS/db"
	"RMS/db/dbHelper"
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
		var req models.RegisterRequestDB
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

		req.CreatedBy = claims.UserID

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
		//pasword
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if hashErr != nil {
			http.Error(w, "password hash error", http.StatusBadRequest)
			return
		}
		req.Password = string(hashedPassword)

		tx, err := db.RM.Begin()
		defer utils.Tx(tx, &err)
		err = dbHelper.SearchEmail(tx, &req, w)
		err = dbHelper.SearchPhone(tx, &req, w)

		err = dbHelper.RegisterUser(tx, &req)

		if err = tx.Commit(); err != nil {
			http.Error(w, "failed to complete registration", http.StatusInternalServerError)
			return
		}
		resp := models.CreateUserResponse{
			UserId:         req.UserId,
			Email:          req.Email,
			Role:           req.Role,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
