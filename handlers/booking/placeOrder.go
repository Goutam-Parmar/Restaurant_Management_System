package booking

import (
	"RMS/models"
	"RMS/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

func PlaceOrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ✅ Extract JWT claims
		claims, err := utils.ExtractAuthClaims(r.Header.Get("Authorization"))
		if err != nil {
			log.Println("unauthorized:", err)
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		userID := claims.UserID

		// ✅ Decode request body
		var req models.PlaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("invalid request body:", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// ✅ Validate required fields
		if req.RestaurantID == 0 || req.AddressID == 0 || len(req.Items) == 0 {
			http.Error(w, "missing required fields", http.StatusBadRequest)
			return
		}

		// ✅ Check restaurant is active
		var isActive bool
		err = db.QueryRow(`SELECT is_active FROM restaurants WHERE id = $1`, req.RestaurantID).Scan(&isActive)
		if err != nil || !isActive {
			http.Error(w, "invalid or inactive restaurant", http.StatusBadRequest)
			return
		}

		// ✅ Check address belongs to user
		var addressExists bool
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM addresses WHERE id = $1 AND user_id = $2)`,
			req.AddressID, userID).Scan(&addressExists)
		if err != nil || !addressExists {
			http.Error(w, "invalid address", http.StatusBadRequest)
			return
		}

		// ✅ Start transaction
		tx, err := db.Begin()
		if err != nil {
			log.Println("failed to begin transaction:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback() // always rollback unless commit succeeds

		// ✅ Insert order with initial total_amount=0
		var orderID int64
		err = tx.QueryRow(`
			INSERT INTO orders (user_id, restaurant_id, address_id, status, payment_status, total_amount, payment_method)
			VALUES ($1, $2, $3, 'placed', 'pending', 0, $4)
			RETURNING id
		`,
			userID, req.RestaurantID, req.AddressID, req.PaymentMethod,
		).Scan(&orderID)
		if err != nil {
			log.Println("failed to create order:", err)
			http.Error(w, "failed to create order", http.StatusInternalServerError)
			return
		}

		// ✅ Insert order items & calculate total
		var totalAmount float64
		menuSeen := make(map[int64]bool) // prevent duplicate menu_id

		for _, item := range req.Items {
			if item.Quantity <= 0 {
				err = fmt.Errorf("invalid quantity for menu_id %d", item.MenuID)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if menuSeen[item.MenuID] {
				err = fmt.Errorf("duplicate menu_id %d", item.MenuID)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			menuSeen[item.MenuID] = true

			var price float64
			err = tx.QueryRow(`
				SELECT price FROM menus
				WHERE id = $1 AND restaurant_id = $2 AND is_available = true
			`, item.MenuID, req.RestaurantID).Scan(&price)
			if err != nil {
				log.Println("invalid menu item:", err)
				http.Error(w, fmt.Sprintf("invalid or unavailable menu_id: %d", item.MenuID), http.StatusBadRequest)
				return
			}

			_, err = tx.Exec(`
				INSERT INTO order_items (order_id, menu_id, quantity, price)
				VALUES ($1, $2, $3, $4)
			`, orderID, item.MenuID, item.Quantity, price)
			if err != nil {
				log.Println("failed to insert order item:", err)
				http.Error(w, "failed to add order item", http.StatusInternalServerError)
				return
			}

			totalAmount += price * float64(item.Quantity)
		}

		// ✅ Round total_amount to 2 decimals if needed
		totalAmount = math.Round(totalAmount*100) / 100

		// ✅ Update order with final total_amount
		_, err = tx.Exec(`UPDATE orders SET total_amount = $1 WHERE id = $2`, totalAmount, orderID)
		if err != nil {
			log.Println("failed to update total amount:", err)
			http.Error(w, "failed to update order amount", http.StatusInternalServerError)
			return
		}

		// ✅ Commit transaction
		if err := tx.Commit(); err != nil {
			log.Println("failed to commit transaction:", err)
			http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
			return
		}

		// ✅ Send response
		resp := models.PlaceOrderResponse{
			Message:        "Order placed successfully",
			OrderID:        orderID,
			TotalAmount:    totalAmount,
			Status:         "placed",
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
