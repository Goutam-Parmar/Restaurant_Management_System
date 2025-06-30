package auth

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func ListAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		role := strings.ToLower(claims.Role)
		if role != "admin" && role != "sub-admin" {
			http.Error(w, "forbidden: only admin or sub-admin can view users", http.StatusForbidden)
			return
		}

		query := `
			SELECT u.id, u.name, a.city
            FROM users u
            JOIN user_roles ur ON u.id = ur.user_id
            LEFT JOIN addresses a ON u.id = a.user_id AND a.is_primary = true
            WHERE ur.role = 'user'
		`
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var users []models.UserBrief2
		for rows.Next() {
			var user models.UserBrief2
			if err := rows.Scan(&user.ID, &user.Name, &user.City); err != nil {
				continue
			}
			users = append(users, user)
		}
		resp := models.GetAllUsersResponse{
			Message:        "User list fetched successfully",
			Users2:         users,
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
