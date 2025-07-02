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

func ListAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		role := strings.ToLower(claims.Role)
		if role == "sub-admin" {
			http.Error(w, "forbidden: only admin or sub-admin can view users", http.StatusForbidden)
			return
		}
		// helper fuction to get list of users
		users, err := dbHelper.GetListOfUsers(w)
		if err != nil {
			http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		}

		resp := models.GetAllUsersResponse{
			Users2:         users,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
