// Точка входа сервиса order: HTTP для клиентов, расчёт цены через gRPC у сервиса pricing.
//
// Ученик: после реализации client.Dial этот main заработает с запущенным pricing.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"golang_grpc_class/services/order/client"
	"golang_grpc_class/services/order/handler"
)

func main() {
	// Адрес gRPC-сервиса ценообразования (должен совпадать с тем, куда слушает pricing).
	pricingAddr := os.Getenv("PRICING_GRPC_ADDR")
	if pricingAddr == "" {
		pricingAddr = "localhost:50051"
	}

	// Dial с таймаутом: при старте pricing должен быть доступен, иначе процесс завершится.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pc, closeConn, err := client.Dial(ctx, pricingAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = closeConn() }()

	// По умолчанию 8081, чтобы не пересекаться с монолитом на 8080.
	addr := ":8081"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	h := handler.New(pc)
	log.Printf("order service listening on %s (pricing gRPC: %s)", addr, pricingAddr)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal(err)
	}
}
