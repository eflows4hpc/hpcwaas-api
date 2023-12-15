package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

// SecureRandomBytes returns the requested number of bytes using crypto/rand
func SecureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

// SecureRandomSecret returns Base64-encoded random string of the specified leangth
func SecureRandomSecret(length int) string {
	randomBytes := SecureRandomBytes(length)
	h := sha256.New()
	h.Write(randomBytes)
	randomString := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return randomString
}
