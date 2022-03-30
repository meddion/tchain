package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigningWithSecretKey(t *testing.T) {
	signer, err := NewSigner()
	assert.NoError(t, err, "on generating a secret key")

	t.Run("signing", func(t *testing.T) {
		msg := []byte("How are you doing?")
		sig, err := signer.Sign(msg)
		assert.NoError(t, err, "on signing a message")

		assert.True(t, sig.Verify(msg), "on verifying a message")
	})

	t.Run("signing_with_hash", func(t *testing.T) {
		msg := []byte("this message should be signed")

		hashedMsg, err := Hash256(msg)
		assert.NoError(t, err, "on hashing a message")

		sig, err := signer.Sign(hashedMsg[:])
		assert.NoError(t, err, "on signing a message")

		assert.True(t, sig.Verify(hashedMsg[:]), "on verifying a message")

	})

	t.Run("signing_wrong_msg_error", func(t *testing.T) {
		msg := []byte("0x000002")
		sig, err := signer.Sign(msg)
		assert.NoError(t, err, "on signing a message")

		assert.False(t, sig.Verify([]byte("0x000001")), "on verifying a message")
	})
}
