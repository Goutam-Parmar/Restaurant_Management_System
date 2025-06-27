package auth

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetAllSubAdmins(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil || claims.Role != "admin" {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT 
                     u.id, 
                     u.name, 
                     u.email, 
                     a.city
          FROM users u
    JOIN user_roles ur ON u.id = ur.user_id
   JOIN addresses a ON u.id = a.user_id AND a.is_primary = true
        WHERE ur.role = 'sub-admin';
		`)
		if err != nil {
			log.Println("❌ Failed to query subadmins:", err)
			http.Error(w, "Failed to fetch subadmins", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.UserBrief
		for rows.Next() {
			var u models.UserBrief
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.City); err != nil {
				log.Println("❌ Row scan error:", err)
				continue
			}
			users = append(users, u)
		}

		// ✅ Prepare response
		resp := models.GetAllSubAdminsResponse{
			Message: "Sub-admins fetched successfully",
			Users:   users,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
