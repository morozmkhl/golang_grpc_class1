// Пакет main — точка входа монолита: один процесс принимает HTTP и сам считает скидки.
package main

import (
	"log"
	"net/http"
	"os"

	"golang_grpc_class/monolith/handler"
	"golang_grpc_class/monolith/service"
)

func main() {
	// Адрес по умолчанию совпадает с README; PORT может быть "8080" или ":8080".
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	svc := service.NewOrderService()
	h := handler.New(svc)

	log.Printf("monolith listening on %s", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal(err)
	}
}
