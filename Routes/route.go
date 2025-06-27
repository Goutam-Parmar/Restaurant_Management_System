package Routes

import (
	"RMS/handlers/Restaurant"
	"RMS/handlers/auth"
	"RMS/utils"
	"database/sql"
	"github.com/gorilla/mux"
)

func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth/register", auth.Register(db)).Methods("POST")
	router.HandleFunc("/auth/login", auth.Login(db)).Methods("POST")

	router.HandleFunc("/auth/refresh", auth.RefreshToken()).Methods("GET")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(utils.JWTMiddleware())
	protected.HandleFunc("/getList/users", auth.ListAllUsers(db)).Methods("GET")
	protected.HandleFunc("/user/{user_id}/address", auth.AddSingleAddress(db)).Methods("POST")

	protected.HandleFunc("/restaurants", restaurant.GetAllRestaurants(db)).Methods("GET")
	protected.HandleFunc("/restaurants/{restaurant_id}/menus", restaurant.GetMenusByRestaurantID(db)).Methods("GET")
	protected.HandleFunc("/restaurants/{restaurant_id}/distance", restaurant.GetDistanceToRestaurant(db)).Methods("GET")

	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(utils.RequireRole("admin"))

	admin.HandleFunc("/restaurants/createNewRestaurant", restaurant.CreateRestaurant(db)).Methods("POST")
	admin.HandleFunc("/restaurants/addMenu/{restaurant_id}", restaurant.AddMenu(db)).Methods("POST")
	admin.HandleFunc("/createSubAdmin", auth.CreateSubAdmin(db)).Methods("POST")
	admin.HandleFunc("/getAllSubAdmins", auth.GetAllSubAdmins(db)).Methods("GET")

	subadmin := protected.PathPrefix("/subadmin").Subrouter()
	subadmin.Use(utils.RequireRole("sub-admin"))
	protected.HandleFunc("/subadmin/create-restaurant", restaurant.CreateRestaurantBySubAdmin(db)).Methods("POST")
	subadmin.HandleFunc("/add-menu/{restaurant_id}", restaurant.AddMenuBySubadmin(db)).Methods("POST")
	return router
}
