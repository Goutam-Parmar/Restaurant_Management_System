package auth

import (
	"RMS/db"
	"RMS/models"
	"RMS/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddSingleAddress() http.HandlerFunc {
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
			http.Error(w, "forbidden: only users can add/change addresses", http.StatusForbidden)
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
		//transaction
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

		if req.IsPrimary {
			_, err = tx.Exec(`UPDATE addresses SET is_primary = false WHERE user_id = $1`, pathUserID)
			if err != nil {
				http.Error(w, "failed to unset other primary addresses", http.StatusInternalServerError)
				return
			}
		}
		var exists bool
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM addresses WHERE user_id = $1 AND label = $2)`, pathUserID, req.Label).Scan(&exists)
		if err != nil {
			http.Error(w, "failed to check existing address", http.StatusInternalServerError)
			return
		}
		if exists {
			_, err = tx.Exec(`
				UPDATE addresses
				SET address_line = $1, city = $2, latitude = $3, longitude = $4, is_primary = $5
				WHERE user_id = $6 AND label = $7
			`,
				req.AddressLine, req.City, req.Latitude, req.Longitude, req.IsPrimary, pathUserID, req.Label,
			)
		} else {
			// Insert
			_, err = tx.Exec(`
				INSERT INTO addresses (user_id, label, address_line, city, latitude, longitude, is_primary)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`,
				pathUserID, req.Label, req.AddressLine, req.City, req.Latitude, req.Longitude, req.IsPrimary,
			)
		}
		if err != nil {
			http.Error(w, "failed to upsert address: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err = tx.Commit(); err != nil {
			http.Error(w, "DB commit error", http.StatusInternalServerError)
			return
		}
		resp := map[string]interface{}{
			"message":          "Address added or updated successfully",
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
