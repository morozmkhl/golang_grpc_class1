package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"

	"golang_grpc_class/pkg/pricingpb"
	"golang_grpc_class/services/order/client"
	"golang_grpc_class/services/order/handler"
	"golang_grpc_class/services/pricing/grpcserver"
	pricingsvc "golang_grpc_class/services/pricing/service"
)

// TestMicroservicesHTTP_SharedCases поднимает pricing (gRPC) и order (HTTP) в памяти и прогоняет SharedCases.
//
// Удалите t.Skip ниже после реализации gRPC-клиента и grpcserver.CalculatePrice (см. solution.md).
func TestMicroservicesHTTP_SharedCases(t *testing.T) {
	t.Skip("реализуйте gRPC по заданию, затем удалите этот Skip — см. solution.md")

	// Слушаем случайный порт — тесты не конфликтуют при параллельном запуске.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	grpcSrv := grpc.NewServer()
	pricingpb.RegisterPricingServiceServer(grpcSrv, &grpcserver.Server{Svc: pricingsvc.NewPricingService()})
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			t.Logf("grpc serve: %v", err)
		}
	}()
	t.Cleanup(func() { grpcSrv.Stop() })

	ctx := context.Background()
	pc, closeConn, err := client.Dial(ctx, lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = closeConn() })

	h := handler.New(pc)
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
