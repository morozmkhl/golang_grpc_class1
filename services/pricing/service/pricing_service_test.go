package service

import (
	"math"
	"testing"
)

// TestPricingService_FinalPrice согласован с монолитом на суммах из SharedCases.
func TestPricingService_FinalPrice(t *testing.T) {
	s := NewPricingService()
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
