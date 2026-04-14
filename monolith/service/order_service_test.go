package service

import (
	"math"
	"testing"
)

// TestOrderService_Discount проверяет ступени скидок и границы (100 и 500).
func TestOrderService_Discount(t *testing.T) {
	s := NewOrderService()
	tests := []struct {
		name   string
		amount float64
		want   float64
	}{
		{"no discount", 50, 0},
		{"10 percent", 150, 0.10},
		{"20 percent", 600, 0.20},
		{"boundary 100 no discount", 100, 0},
		{"just above 100", 100.01, 0.10},
		{"boundary 500 ten percent", 500, 0.10},
		{"just above 500", 500.01, 0.20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.Discount(tt.amount)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("Discount(%v) = %v, want %v", tt.amount, got, tt.want)
			}
		})
	}
}

// TestOrderService_FinalPrice — ключевые суммы 50 / 150 / 600 из README.
func TestOrderService_FinalPrice(t *testing.T) {
	s := NewOrderService()
	tests := []struct {
		amount float64
		want   float64
	}{
		{50, 50},
		{150, 135},
		{600, 480},
	}
	for _, tt := range tests {
		got := s.FinalPrice(tt.amount)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Fatalf("FinalPrice(%v) = %v, want %v", tt.amount, got, tt.want)
		}
	}
}
