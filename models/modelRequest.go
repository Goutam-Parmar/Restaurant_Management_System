package models

type RegisterRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Phone       string  `json:"phone"`
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
	Rating    float64 `json:"rating"`
}
type AddMenuRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	FoodType    string  `json:"food_type"`
	Category    string  `json:"category"`
}
type AddAddressRequest struct {
	Label       string  `json:"label"`
	AddressLine string  `json:"address_line"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	IsPrimary   bool    `json:"is_primary"`
}
type OrderItemRequest struct {
	MenuID   int64 `json:"menu_id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required,gt=0"`
}

type PlaceOrderRequest struct {
	RestaurantID  int64              `json:"restaurant_id" validate:"required"`
	AddressID     int64              `json:"address_id" validate:"required"`
	Items         []OrderItemRequest `json:"items" validate:"required,dive"`
	PaymentMethod string             `json:"payment_method"`
}
