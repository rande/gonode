package vault

import (
	"crypto/hmac"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test hmac usage
func Test_Hmac(t *testing.T) {
	mac := hmac.New(sha256.New, key)
	mac.Write(xLargeMessage)
	macFull := mac.Sum(nil)

	mac = hmac.New(sha256.New, key)

	for _, b := range xLargeMessage {
		mac.Write([]byte{b})
	}

	macChunk := mac.Sum(nil)

	assert.Equal(t, macChunk, macFull)
}
