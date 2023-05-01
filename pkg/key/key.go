package key

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func WrapKey(key string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	randomString := make([]byte, 4)
	for i := range randomBytes {
		randomString[i] = charset[int(randomBytes[i])%len(charset)]
	}
	authKey := fmt.Sprintf("%s:%s", key, string(randomString))
	encodedAuthKey := base64.RawURLEncoding.EncodeToString([]byte(authKey))
	return encodedAuthKey
}

func ExtractKey(inviteKey string) string {
	decodedKey, err := base64.RawURLEncoding.DecodeString(inviteKey)
	if err != nil {
		return ""
	}
	authKey := string(decodedKey)

	return authKey
}
