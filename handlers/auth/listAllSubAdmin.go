package auth

import (
	"RMS/db/dbHelper"
	"RMS/models"
	"RMS/utils"
	"encoding/json"

	"net/http"
)

func GetAllSubAdmins() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil || claims.Role != "admin" {
			http.Error(w, "Unauthorized access, only admin can acess", http.StatusUnauthorized)
			return
		}
		users, err := dbHelper.GetAllSubAdmin(w)
		if err != nil {
			http.Error(w, "Failed to fetch sub-admins", http.StatusInternalServerError)
			return
		}
		resp := models.GetAllUsersResponse{
			Users2: users,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
