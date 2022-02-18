package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyGen(t *testing.T) {
	sk, err := NewSecretKey()
	assert.NoError(t, err, "on generating a secret key")

	pk := sk.PublicKey()
	assert.True(t, pk.IsValid(), "on checking key validity")

	pk = PublicKey{}
	assert.False(t, pk.IsValid(), "on checking key validity")
}

func TestSigningWithSecretKey(t *testing.T) {
	sk, err := NewSecretKey()
	assert.NoError(t, err, "on generating a secret key")

	t.Run("signing", func(t *testing.T) {
		msg := []byte("How are you doing?")
		r, s, err := sk.Sign(msg)
		assert.NoError(t, err, "on signing a message")

		pk := sk.PublicKey()
		assert.True(t, Verify(pk, msg, r, s), "on verifying a message")
	})

	t.Run("signing_with_hash", func(t *testing.T) {
		msg := []byte("this message should be signed")

		hashedMsg, err := Hash256(msg)
		assert.NoError(t, err, "on hashing a message")

		r, s, err := sk.Sign(hashedMsg[:])
		assert.NoError(t, err, "on signing a message")

		pk := sk.PublicKey()
		assert.True(t, Verify(pk, msg, r, s), "on verifying a message")

	})

	t.Run("signing_error", func(t *testing.T) {
		msg := []byte("0x000002")
		r, s, err := sk.Sign(msg)
		assert.NoError(t, err, "on signing a message")

		pk := sk.PublicKey()
		assert.False(t, Verify(pk, []byte("0x000001"), r, s), "on verifying a message")
	})
}
