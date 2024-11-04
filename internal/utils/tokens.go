package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateToken creates a secure random token for CSRF
func GenerateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
