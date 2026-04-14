// Пакет handler — HTTP-слой монолита: разбор JSON, вызов доменного сервиса, ответ клиенту.
package handler

import (
	"encoding/json"
	"net/http"

	"golang_grpc_class/monolith/service"
)

// OrderHTTP реализует HTTP API монолита для создания заказов.
type OrderHTTP struct {
	svc *service.OrderService
}

// New собирает mux с маршрутом POST /order и возвращает готовый http.Handler.
func New(svc *service.OrderService) http.Handler {
	h := &OrderHTTP{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /order", h.handleCreateOrder)
	return mux
}

// Тело запроса POST /order (совпадает с микросервисом order).
type createOrderRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

// Тело успешного ответа: id заказа и цена после скидки.
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
	order, err := h.svc.CreateOrder(req.UserID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Итоговая цена считается тем же сервисом, что и скидки — без отдельного gRPC.
	resp := createOrderResponse{
		OrderID:    order.ID,
		FinalPrice: h.svc.FinalPrice(req.Amount),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}
