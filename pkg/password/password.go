package password

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate a password of len
func Generate(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	// Base 64 can be longer than len
	return str[:len]
}
