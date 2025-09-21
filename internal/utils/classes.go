package utils

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func GenerateRandomDigits(n int) string {
	// Generate random bytes
	b := make([]byte, n)
	_, _ = rand.Read(b)

	// Encode into base32 (letters + digits), cut to required length
	code := strings.ToLower(base32.StdEncoding.EncodeToString(b))
	return code[:n]
}
