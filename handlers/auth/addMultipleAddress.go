package auth

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddSingleAddress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		params := mux.Vars(r)
		pathUserIDStr := params["user_id"]
		pathUserID, err := strconv.ParseInt(pathUserIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user_id in path", http.StatusBadRequest)
			return
		}

		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if strings.ToLower(claims.Role) != "user" {
			http.Error(w, "forbidden: only users can add address", http.StatusForbidden)
			return
		}
		if claims.UserID != pathUserID {
			http.Error(w, "forbidden: user ID mismatch", http.StatusForbidden)
			return
		}

		var req models.AddAddressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = db.Exec(query, pathUserID, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude, req.IsPrimary)
		if err != nil {
			http.Error(w, "failed to add address", http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"message":          "Address added successfully",
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
