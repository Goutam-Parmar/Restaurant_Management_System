package models

type RegisterRequestDB struct {
	UserId      int64   `json:"user_id"`
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
	CreatedBy   int64   `json:"created_By"`
}

type LoginRequest struct {
	UserId   int64  `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
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
	UserId      int64   `json:"user_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	FoodType    string  `json:"food_type"`
	Category    string  `json:"category"`
}
type AddNewAddressRequest struct {
	ID          int64   `json:"id"`
	UserId      int64   `json:"user_id"`
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
