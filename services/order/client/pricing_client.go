// Пакет client — gRPC-клиент к сервису pricing: установка соединения и вызов CalculatePrice.
//
// Ученик: реализуйте Dial и FinalPrice (см. solution.md или материалы курса).
package client

import (
	"context"
	"fmt"
)

// PricingClient оборачивает сгенерированный gRPC-клиент в узкий API приложения.
type PricingClient struct{}

// Dial устанавливает соединение с pricing по addr (например "localhost:50051").
// Второе возвращаемое значение — функция закрытия соединения; её нужно вызвать при остановке.
func Dial(ctx context.Context, addr string) (*PricingClient, func() error, error) {
	_ = ctx
	_ = addr
	return nil, func() error { return nil }, fmt.Errorf("TODO: реализуйте Dial (grpc.DialContext, NewPricingServiceClient)")
}

// FinalPrice вызывает RPC CalculatePrice и возвращает итоговую цену после скидок.
func (c *PricingClient) FinalPrice(ctx context.Context, amount float64) (float64, error) {
	_ = c
	_ = ctx
	_ = amount
	return 0, fmt.Errorf("TODO: реализуйте FinalPrice (вызов CalculatePrice)")
}
