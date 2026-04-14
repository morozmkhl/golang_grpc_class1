// gRPC-сервер pricing: единственная ответственность — расчёт итоговой цены по правилам скидок.
//
// Ученик: реализуйте точку входа — net.Listen, grpc.NewServer, RegisterPricingServiceServer,
// опционально reflection, Serve (см. solution.md).
package main

import (
	"log"
)

func main() {
	log.Fatal("TODO: реализуйте gRPC-сервер pricing (см. solution.md)")
}
