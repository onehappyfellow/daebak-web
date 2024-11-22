package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("random bytes: %w", err)
	}
	return b, nil
}

func RandomString(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", fmt.Errorf("random string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
