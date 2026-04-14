// HTTP-слой сервиса order: тот же контракт, что у монолита, но final_price приходит из gRPC.
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang_grpc_class/services/order/client"
)

// OrderHTTP обрабатывает запросы к API заказов; ценообразование делегируется микросервису pricing.
type OrderHTTP struct {
	pricing *client.PricingClient
}

// New собирает mux с POST /order.
func New(pricing *client.PricingClient) http.Handler {
	h := &OrderHTTP{pricing: pricing}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /order", h.handleCreateOrder)
	return mux
}

type createOrderRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type createOrderResponse struct {
	OrderID    string  `json:"order_id"`
	FinalPrice float64 `json:"final_price"`
}

func (h *OrderHTTP) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.Amount < 0 {
		http.Error(w, "amount must be non-negative", http.StatusBadRequest)
		return
	}

	// Отдельный таймаут на gRPC: медленный pricing не должен держать запрос бесконечно.
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	final, err := h.pricing.FinalPrice(ctx, req.Amount)
	if err != nil {
		// 502: ошибка не в теле запроса клиента, а у зависимости (pricing).
		http.Error(w, "pricing unavailable", http.StatusBadGateway)
		return
	}

	// Генерация order_id на границе HTTP — идентификатор заказа принадлежит сервису order.
	orderID, err := randomOrderID()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := createOrderResponse{
		OrderID:    orderID,
		FinalPrice: final,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}
