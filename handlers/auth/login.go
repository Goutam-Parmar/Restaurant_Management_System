package auth

import (
	"RMS/db/dbHelper"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func Login() http.HandlerFunc {
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
		err := dbHelper.CheckLoginCredentials(&req, w)
		accessToken, err := utils.GenerateAccessToken(req.UserId, req.Email, req.Role)
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(req.UserId, req.Email, req.Role)
		if err != nil {
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}
		resp := models.LoginResponse{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	}
}
