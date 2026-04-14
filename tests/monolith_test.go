package tests

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang_grpc_class/monolith/handler"
	"golang_grpc_class/monolith/service"
)

// TestMonolithHTTP_SharedCases проверяет монолитный HTTP POST /order на всех сценариях из SharedCases.
func TestMonolithHTTP_SharedCases(t *testing.T) {
	h := handler.New(service.NewOrderService())
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	for _, tc := range SharedCases {
		t.Run(tc.Name, func(t *testing.T) {
			body := map[string]any{"user_id": "student", "amount": tc.Amount}
			buf, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.Post(srv.URL+"/order", "application/json", bytes.NewReader(buf))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status %d", resp.StatusCode)
			}
			var out struct {
				OrderID    string  `json:"order_id"`
				FinalPrice float64 `json:"final_price"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				t.Fatal(err)
			}
			if out.OrderID == "" {
				t.Fatal("expected order_id")
			}
			if math.Abs(out.FinalPrice-tc.Want) > 1e-9 {
				t.Fatalf("final_price = %v, want %v", out.FinalPrice, tc.Want)
			}
		})
	}
}
