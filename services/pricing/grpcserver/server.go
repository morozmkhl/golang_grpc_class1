// Адаптер gRPC: перевод protobuf-запросов в вызовы доменного PricingService.
//
// Ученик: реализуйте метод CalculatePrice (см. solution.md).
package grpcserver

import (
	"golang_grpc_class/pkg/pricingpb"
	"golang_grpc_class/services/pricing/service"
)

// Server реализует контракт PricingService из proto/pricing.proto.
// Ученик: добавьте метод CalculatePrice(ctx, *CalculatePriceRequest) (*CalculatePriceResponse, error).
type Server struct {
	pricingpb.UnimplementedPricingServiceServer
	Svc *service.PricingService
}
