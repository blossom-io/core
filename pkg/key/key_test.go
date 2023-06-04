package key

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapKey(t *testing.T) {
	state := "104717035"
	result := WrapKey(state)

	// Ensure the result is not empty
	assert.NotEmpty(t, result)

	// Ensure the result is a valid base64-encoded string
	decoded, err := base64.RawURLEncoding.DecodeString(result)
	assert.NoError(t, err)
	assert.NotEmpty(t, decoded)

	// Ensure the decoded result contains the state and a 4-char random string separated by a colon
	expected := fmt.Sprintf("%s:%s", state, decoded[len(state)+1:])
	assert.Equal(t, expected, string(decoded))
}

func BenchmarkWrapKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WrapKey("olyashaa")
	}
}

func TestExtractKey(t *testing.T) {
	inviteKey := "104717035:6m1w"
	encodedAuthKey := "MTA0NzE3MDM1OjZtMXc"

	result := ExtractKey(encodedAuthKey)

	// Ensure the result matches the expected state
	assert.Equal(t, inviteKey, result)

	// // Ensure the result is empty if the input is not a valid base64-encoded string
	result = ExtractKey("not-a-valid-base64-string")
	assert.Empty(t, result)
}
