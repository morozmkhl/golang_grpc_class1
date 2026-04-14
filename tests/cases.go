// Пакет tests — интеграционные тесты: сравнивают HTTP ответы монолита и связки order→pricing (gRPC).
// Общие сценарии вынесены в SharedCases, чтобы оба пути давали одинаковый final_price.
package tests

// SharedCases — эталонные суммы и ожидаемый final_price (должны совпадать у монолита и микросервисов).
var SharedCases = []struct {
	Name   string
	Amount float64
	Want   float64
}{
	{"no discount", 50, 50},
	{"10 percent discount", 150, 135},
	{"20 percent discount", 600, 480},
}
