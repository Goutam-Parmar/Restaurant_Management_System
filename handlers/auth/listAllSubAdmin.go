package auth

import (
	"RMS/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetAllSubAdmins(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT u.id, u.name, u.email, COALESCE(a.city, '') AS city
			FROM users u
			JOIN user_roles ur ON u.id = ur.user_id
			LEFT JOIN addresses a ON u.id = a.user_id AND a.is_primary = true
			WHERE ur.role = 'sub-admin'
		`)
		if err != nil {
			log.Println("faail to query subadmins:", err)
			http.Error(w, "fail to fetch subadmins", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.UserBrief

		for rows.Next() {
			var u models.UserBrief
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.City); err != nil {
				log.Println("row inseert fail error:", err)
				continue
			}
			users = append(users, u)
		}

		resp := models.GetAllSubAdminsResponse{
			Message: "sub-admins fetched successfully",
			Users:   users,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
