package models

type RegisteredUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
type RegisterResponse struct {
	Message        string         `json:"message"`
	User           RegisteredUser `json:"user"`
	AccessToken    string         `json:"access_token"`
	RefreshToken   string         `json:"refresh_token"`
	ResponseTimeMs float64        `json:"response_time_ms"`
}
type LoginUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
type LoginResponse struct {
	Message        string    `json:"message"`
	User           LoginUser `json:"user"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ResponseTimeMs float64   `json:"response_time_ms"`
}
type RefreshResponse struct {
	Message        string  `json:"message"`
	AccessToken    string  `json:"access_token"`
	RefreshToken   string  `json:"refresh_token"`
	ResponseTimeMs float64 `json:"response_time_ms"`
}

type CreatedRestaurant struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rating    int     `json:"rating"`
	IsActive  bool    `json:"is_active"`
	CreatedBy int64   `json:"created_by"`
}
type CreateRestaurantResponse struct {
	Message        string            `json:"message"`
	Restaurant     CreatedRestaurant `json:"restaurant"`
	ResponseTimeMs float64           `json:"response_time_ms"`
}
type AddedMenu struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"is_available"`
	FoodType     string  `json:"food_type"`
	Category     string  `json:"category"`
	RestaurantID int64   `json:"restaurant_id"`
	CreatedBy    int64   `json:"created_by"`
}
type AddMenuResponse struct {
	Message        string    `json:"message"`
	Menu           AddedMenu `json:"menu"`
	ResponseTimeMs float64   `json:"response_time_ms"`
}
type SlimRestaurantResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Rating   int    `json:"rating"`
	IsActive bool   `json:"is_active"`
}
type AllRestaurantsSlimResponse struct {
	Message        string                   `json:"message"`
	Restaurants    []SlimRestaurantResponse `json:"restaurants"`
	ResponseTimeMs float64                  `json:"response_time_ms"`
}
type MenuItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	FoodType    string  `json:"food_type"`
	Category    string  `json:"category"`
}

type MenuListResponse struct {
	Message        string     `json:"message"`
	RestaurantID   int64      `json:"restaurant_id"`
	RestaurantName string     `json:"restaurant_name"`
	Menus          []MenuItem `json:"menus"`
	ResponseTimeMs float64    `json:"response_time_ms"`
}

type MinimalRestaurant struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type RestaurantDistanceResponse struct {
	Message        string            `json:"message"`
	DistanceKM     float64           `json:"distance_km"`
	Restaurant     MinimalRestaurant `json:"restaurant"`
	ResponseTimeMs float64           `json:"response_time_ms"`
}
type CreateUserResponse struct {
	Message        string  `json:"message"`
	UserID         int64   `json:"user_id"`
	Email          string  `json:"email"`
	ResponseTimeMs float64 `json:"response_time_ms"`
}
type UserBrief struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	City  string `json:"city,omitempty"`
	//CreatedAt string `json:"created_at"`
}

type GetAllSubAdminsResponse struct {
	Message string      `json:"message"`
	Users   []UserBrief `json:"users"`
}

type AddMenuResponseSubAdmin struct {
	Message        string  `json:"message"`
	MenuID         int64   `json:"menu_id"`
	ResponseTimeMs float64 `json:"response_time_ms"`
}
type UserBrief2 struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type GetAllUsersResponse struct {
	Message        string       `json:"message"`
	Users2         []UserBrief2 `json:"users"`
	ResponseTimeMs float64      `json:"response_time_ms"`
}
