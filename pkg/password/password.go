package password

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate a password of len
func Generate(len int) string {
	buff := make([]byte, len)
	_, err := rand.Read(buff)
	if err != nil {
		return "unable-to-create-a-random-password"[:len]
	}
	str := base64.StdEncoding.EncodeToString(buff)
	// Base 64 can be longer than len
	return str[:len]
}
