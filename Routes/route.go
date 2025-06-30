package Routes

import (
	"RMS/handlers/Restaurant"
	"RMS/handlers/auth"
	"database/sql"
	"github.com/gorilla/mux"
)

func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/auth/user/registerByAdmin", auth.RegisterNewUser(db)).Methods("POST")
	router.HandleFunc("/api/v1/auth/user/login", auth.Login(db)).Methods("POST")

	router.HandleFunc("/api/v1/auth/user/refresh", auth.RefreshToken()).Methods("GET")
	router.HandleFunc("/api/v1/getList/users", auth.ListAllUsers(db)).Methods("GET")
	router.HandleFunc("/api/v1/user/{user_id}/address", auth.AddSingleAddress(db)).Methods("POST")

	router.HandleFunc("/api/v1/list/getRestaurantList", restaurant.GetAllRestaurants(db)).Methods("GET")
	router.HandleFunc("/api/v1/getMenuList/{restaurant_id}", restaurant.GetMenusByRestaurantID(db)).Methods("GET")
	router.HandleFunc("/api/v1/restaurants/getdistance/{restaurant_id}", restaurant.GetDistanceToRestaurant(db)).Methods("GET")

	admin := router.PathPrefix("/api/v1/admin").Subrouter()

	admin.HandleFunc("/restaurants/createNewRestaurant", restaurant.CreateRestaurantByAdmin(db)).Methods("POST")
	admin.HandleFunc("/restaurants/addMenu/{restaurant_id}", restaurant.AddMenuByAdmin(db)).Methods("POST")
	admin.HandleFunc("/getAllSubAdmins", auth.GetAllSubAdmins(db)).Methods("GET")

	subadmin := router.PathPrefix("/subadmin").Subrouter()

	subadmin.HandleFunc("/api/v1/create-restaurant", restaurant.CreateRestaurantBySubAdmin(db)).Methods("POST")
	subadmin.HandleFunc("/api/v1/addMenuBySubAdmin/{restaurant_id}", restaurant.AddMenuBySubadmin(db)).Methods("POST")
	subadmin.HandleFunc("/api/v1/auth/user/registerBySubAdmin", auth.CreateUserBySubAdmin(db)).Methods("POST")
	return router
}
