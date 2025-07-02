package auth

import (
	"RMS/db"
	"RMS/db/dbHelper"
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
		var req models.AddNewAddressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		req.UserId = claims.UserID
		tx, err := db.RM.Begin()
		defer utils.Tx(tx, &err)
		err = dbHelper.AddNewAddress(tx, &req, w)
		if err = tx.Commit(); err != nil {
			http.Error(w, "failed to complete registration", http.StatusInternalServerError)
			return
		}
		resp := map[string]interface{}{
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
