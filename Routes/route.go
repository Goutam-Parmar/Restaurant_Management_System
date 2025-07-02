package Routes

import (
	"RMS/handlers/Restaurant"
	"RMS/handlers/auth"
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/auth/user/registerByAdmin", auth.RegisterNewUser()).Methods("POST")
	router.HandleFunc("/api/v1/auth/user/login", auth.Login()).Methods("POST")

	router.HandleFunc("/api/v1/auth/user/refresh", auth.RefreshToken()).Methods("GET")
	router.HandleFunc("/api/v1/getList/users", auth.ListAllUsers()).Methods("GET")
	router.HandleFunc("/api/v1/user/{user_id}/address", auth.AddSingleAddress()).Methods("POST")

	router.HandleFunc("/api/v1/list/getRestaurantList", restaurant.GetAllRestaurants()).Methods("GET")
	router.HandleFunc("/api/v1/getMenuList/{restaurant_id}", restaurant.GetMenusByRestaurantID()).Methods("GET")
	router.HandleFunc("/api/v1/restaurants/getdistance/{restaurant_id}", restaurant.GetDistanceToRestaurant()).Methods("GET")

	admin := router.PathPrefix("/api/v1/admin").Subrouter()

	admin.HandleFunc("/restaurants/createNewRestaurant", restaurant.CreateRestaurantByAdmin()).Methods("POST")
	admin.HandleFunc("/restaurants/addMenu/{restaurant_id}", restaurant.AddMenuByAdmin()).Methods("POST")
	admin.HandleFunc("/getAllSubAdmins", auth.GetAllSubAdmins()).Methods("GET")

	subadmin := router.PathPrefix("/subadmin").Subrouter()

	subadmin.HandleFunc("/api/v1/create-restaurant", restaurant.CreateRestaurantBySubAdmin()).Methods("POST")
	subadmin.HandleFunc("/api/v1/addMenuBySubAdmin/{restaurant_id}", restaurant.AddMenuBySubadmin()).Methods("POST")
	subadmin.HandleFunc("/api/v1/auth/user/registerBySubAdmin", auth.CreateUserBySubAdmin()).Methods("POST")
	return router
}
