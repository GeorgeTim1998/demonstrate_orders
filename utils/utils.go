package utils

import (
	"math/rand"

	"github.com/google/uuid"
)

func GenerateRandomUID() string {
	return uuid.New().String()
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func GenerateRandomPhone() string {
	return GenerateRandomString(10)
}

func GenerateRandomEmail() string {
	return GenerateRandomString(8) + "@example.com"
}
