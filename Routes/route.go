package Routes

import (
	"RMS/handlers/Restaurant"
	"RMS/handlers/auth"
	"database/sql"
	"github.com/gorilla/mux"
)

func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth/register", auth.Register(db)).Methods("POST")
	router.HandleFunc("/auth/login", auth.Login(db)).Methods("POST")

	router.HandleFunc("/auth/refresh", auth.RefreshToken()).Methods("GET")
	router.HandleFunc("/getList/users", auth.ListAllUsers(db)).Methods("GET")
	router.HandleFunc("/user/{user_id}/address", auth.AddSingleAddress(db)).Methods("POST")

	router.HandleFunc("/restaurants", restaurant.GetAllRestaurants(db)).Methods("GET")
	router.HandleFunc("/restaurants/{restaurant_id}/menus", restaurant.GetMenusByRestaurantID(db)).Methods("GET")
	router.HandleFunc("/restaurants/{restaurant_id}/distance", restaurant.GetDistanceToRestaurant(db)).Methods("GET")

	admin := router.PathPrefix("/admin").Subrouter()

	admin.HandleFunc("/restaurants/createNewRestaurant", restaurant.CreateRestaurant(db)).Methods("POST")
	admin.HandleFunc("/restaurants/addMenu/{restaurant_id}", restaurant.AddMenu(db)).Methods("POST")
	admin.HandleFunc("/createSubAdmin", auth.CreateSubAdmin(db)).Methods("POST")
	admin.HandleFunc("/getAllSubAdmins", auth.GetAllSubAdmins(db)).Methods("GET")

	subadmin := router.PathPrefix("/subadmin").Subrouter()

	subadmin.HandleFunc("/create-restaurant", restaurant.CreateRestaurantBySubAdmin(db)).Methods("POST")
	subadmin.HandleFunc("/add-menu/{restaurant_id}", restaurant.AddMenuBySubadmin(db)).Methods("POST")
	return router
}
