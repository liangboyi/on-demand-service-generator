package adapter

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func newRequirepass() (string, error) {
	randomBytes := make([]byte, 20)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Printf("Error generating random bytes, %v", err)
		return "", err
	}
	return encodeRequirepass(randomBytes), nil
}

func encodeRequirepass(password []byte) string {
	encodedBytes := make([]byte, base64.StdEncoding.EncodedLen(len(password)))
	base64.StdEncoding.Encode(encodedBytes, password)
	return string(encodedBytes)
}
