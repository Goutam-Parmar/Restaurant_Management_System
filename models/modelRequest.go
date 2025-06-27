package models

type RegisterRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Role        string  `json:"role"`
	Label       string  `json:"label"`
	AddressLine string  `json:"address_line"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateRestaurantRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rating    int     `json:"rating"`
}
type AddMenuRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	FoodType    string  `json:"food_type"`
	Category    string  `json:"category"`
}
