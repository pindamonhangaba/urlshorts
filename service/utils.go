package service

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	// DefaultCodeLength is the default length of a generated code
	DefaultCodeLength = 8

	// CharSet defines the characters used for generating random codes
	CharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateRandomCode generates a random code of the given length
func GenerateRandomCode(length int) (string, error) {
	if length <= 0 {
		length = DefaultCodeLength
	}

	var result strings.Builder
	charSetLength := big.NewInt(int64(len(CharSet)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			return "", err
		}
		result.WriteByte(CharSet[randomIndex.Int64()])
	}

	return result.String(), nil
}

// ValidateAPIKey validates the API key against the expected key
func ValidateAPIKey(providedKey, expectedKey string) bool {
	return providedKey == expectedKey
}

// BuildShortURL builds a complete short URL from the base URL, code, and pretty name
func BuildShortURL(baseURL, code, prettyName string) string {
	if prettyName != "" {
		return baseURL + "/" + code + "/" + prettyName
	}
	return baseURL + "/" + code
}
