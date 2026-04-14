// Пакет service — доменная логика монолита: заказы, скидки и итоговая цена в одном месте.
package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Order — заказ в предметной модели (идентификатор, пользователь, исходная сумма).
type Order struct {
	ID     string
	UserID string
	Amount float64
}

// OrderService инкапсулирует правила скидок и создание заказа.
type OrderService struct{}

// NewOrderService создаёт экземпляр OrderService.
func NewOrderService() *OrderService {
	return &OrderService{}
}

// Discount возвращает долю скидки от 0 до 1 для заданной суммы (правила — в README).
func (s *OrderService) Discount(amount float64) float64 {
	switch {
	case amount > 500:
		return 0.20
	case amount > 100:
		return 0.10
	default:
		return 0
	}
}

// FinalPrice применяет скидку к сумме и возвращает цену к оплате.
func (s *OrderService) FinalPrice(amount float64) float64 {
	d := s.Discount(amount)
	return amount * (1 - d)
}

// CreateOrder проверяет входные данные и создаёт заказ со сгенерированным ID.
func (s *OrderService) CreateOrder(userID string, amount float64) (Order, error) {
	if userID == "" {
		return Order{}, fmt.Errorf("user_id is required")
	}
	if amount < 0 {
		return Order{}, fmt.Errorf("amount must be non-negative")
	}
	id, err := generateID()
	if err != nil {
		return Order{}, err
	}
	return Order{
		ID:     id,
		UserID: userID,
		Amount: amount,
	}, nil
}

// generateID формирует случайный 128-битный идентификатор в виде hex-строки.
func generateID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
