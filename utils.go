package porgs

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomBase64String generates a random base64 string of given length
func RandomBase64String(length int) (string, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
