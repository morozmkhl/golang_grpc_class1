package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang_grpc_class/monolith/service"
)

// TestPOSTOrder — успешный POST /order: final_price соответствует скидке 10% на сумму 150.
func TestPOSTOrder(t *testing.T) {
	svc := service.NewOrderService()
	h := New(svc)

	body := map[string]any{"user_id": "u1", "amount": 150.0}
	buf, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		OrderID    string  `json:"order_id"`
		FinalPrice float64 `json:"final_price"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.OrderID == "" {
		t.Fatal("expected order_id")
	}
	if resp.FinalPrice != 135 {
		t.Fatalf("final_price = %v, want 135", resp.FinalPrice)
	}
}

// TestPOSTOrder_InvalidJSON — битый JSON даёт 400.
func TestPOSTOrder_InvalidJSON(t *testing.T) {
	svc := service.NewOrderService()
	h := New(svc)
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
}
