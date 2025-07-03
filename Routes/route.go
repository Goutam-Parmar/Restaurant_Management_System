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
	router.HandleFunc("/api/v1/users/getList", auth.ListAllUsers()).Methods("GET")
	router.HandleFunc("/api/v1/user/{user_id}/addAddress", auth.AddSingleAddress()).Methods("POST")

	router.HandleFunc("/api/v1/restaurants/list", restaurant.GetAllRestaurants()).Methods("GET")
	router.HandleFunc("/api/v1/restaurant/{restaurant_id}/menuList", restaurant.GetMenusByRestaurantID()).Methods("GET")
	router.HandleFunc("/api/v1/restaurant/{restaurant_id}/getDistance", restaurant.GetDistanceToRestaurant()).Methods("GET")

	admin := router.PathPrefix("/api/v1/admin").Subrouter()

	admin.HandleFunc("/restaurants/createNewRestaurant", restaurant.CreateRestaurantByAdmin()).Methods("POST")
	admin.HandleFunc("/restaurants/addMenu/{restaurant_id}", restaurant.RegisterNewUserByAdmin()).Methods("POST")
	admin.HandleFunc("/getAllSubAdmins", auth.GetAllSubAdmins()).Methods("GET")

	subadmin := router.PathPrefix("api/v1/subAdmin").Subrouter()

	subadmin.HandleFunc("/restaurant/createRestaurantBySubAdmin", restaurant.CreateRestaurantBySubAdmin()).Methods("POST")
	subadmin.HandleFunc("/restaurant/{restaurant_id}/addMenuBySubAdmin", restaurant.CreateUserBySubAdmin()).Methods("POST")
	subadmin.HandleFunc("/userRegisterBySubAdmin", auth.RegisterNewUser()).Methods("POST")
	return router
}
