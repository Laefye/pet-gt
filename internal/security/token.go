package security

import "crypto/rand"

const tokenLength = 32

func GenerateToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, tokenLength)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate token: " + err.Error())
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
