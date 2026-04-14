// Доменная логика только ценообразования — та же модель скидок, что и в монолите.
package service

// PricingService содержит правила скидок; дублирует математику OrderService монолита.
type PricingService struct{}

// NewPricingService создаёт сервис ценообразования.
func NewPricingService() *PricingService {
	return &PricingService{}
}

// Discount возвращает долю скидки от 0 до 1 для суммы заказа.
func (s *PricingService) Discount(amount float64) float64 {
	switch {
	case amount > 500:
		return 0.20
	case amount > 100:
		return 0.10
	default:
		return 0
	}
}

// FinalPrice возвращает сумму к оплате после применения скидки.
func (s *PricingService) FinalPrice(amount float64) float64 {
	d := s.Discount(amount)
	return amount * (1 - d)
}
