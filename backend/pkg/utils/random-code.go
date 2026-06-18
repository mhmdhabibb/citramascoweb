package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomNumberCode generates a random numeric string of the specified length.
func GenerateRandomNumberCode(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}
	return string(result), nil
}
