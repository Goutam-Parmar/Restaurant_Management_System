package auth

import (
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, claims, err := utils.ParseToken(tokenStr)
		if err != nil || !token.Valid || !utils.IsTokenType(claims, "refresh") {
			http.Error(w, "invalid or expired refresh token or ", http.StatusUnauthorized)
			return
		}

		authClaims, err := utils.ExtractAuthClaims(authHeader)
		if err != nil {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		accessToken, err := utils.GenerateAccessToken(authClaims.UserID, authClaims.Email, authClaims.Role)
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(authClaims.UserID, authClaims.Email, authClaims.Role)
		if err != nil {
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		resp := models.RefreshResponse{
			Message:        "Token refreshed successfully",
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
