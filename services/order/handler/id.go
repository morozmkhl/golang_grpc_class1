package handler

import (
	"crypto/rand"
	"encoding/hex"
)

// randomOrderID возвращает случайный hex-ID длины 32 символа (16 байт энтропии).
func randomOrderID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
