package api

import (
	"testing"

	"github.com/meddion/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTransactionValidation(t *testing.T) {
	sk, err := crypto.NewSecretKey()
	assert.NoError(t, err)

	msg := []byte(`message to be hashed`)
	hashed, err := crypto.Hash(msg)
	assert.NoError(t, err, "on hashing a message")

	r, s, err := sk.Sign(hashed[:])
	assert.NoError(t, err, "on signing a message")

	testTable := []struct {
		tx     Transaction
		target bool
	}{
		{Transaction{Data: msg, Hash: hashed}, false},
		{Transaction{}, false},
		{Transaction{Data: make([]byte, 0), Hash: hashed}, false},
		{Transaction{PublicKey: sk.PublicKey(), Data: msg, Hash: hashed, R: nil, S: s}, false},
		{Transaction{PublicKey: crypto.PublicKey{}, Data: msg, Hash: hashed, R: r, S: s}, false},
		// Success
		{Transaction{PublicKey: sk.PublicKey(), Data: msg, Hash: hashed, R: r, S: s}, true},
	}

	for i, testCase := range testTable {
		assert.Truef(t, VerifyTransaction(testCase.tx) == testCase.target, "table entry by %d index", i)
	}
}
