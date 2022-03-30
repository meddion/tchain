package api

import (
	"testing"

	"github.com/meddion/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTransactionValidation(t *testing.T) {
	var msg TxData
	copy(msg[:], []byte(`message to be hashed`))

	hashed, err := crypto.Hash256(msg[:])
	assert.NoError(t, err, "on hashing a message")

	signer, err := crypto.NewSigner()
	assert.NoError(t, err, "on creating a signer")
	sig, err := signer.Sign(hashed[:])
	assert.NoError(t, err, "on signing a message")

	testTable := []struct {
		tx     Transaction
		target bool
	}{
		{Transaction{Data: msg, Hash: hashed}, false},
		{Transaction{}, false},
		{Transaction{Data: TxData{}, Hash: hashed}, false},
		{Transaction{Sig: sig, Data: msg, Hash: hashed}, false},
		{Transaction{Sig: sig, Data: msg, Hash: hashed}, false},
		// Success
		{Transaction{Sig: sig, Data: msg, Hash: hashed}, true},
	}

	for i, testCase := range testTable {
		assert.Truef(t, VerifyTransaction(testCase.tx) == testCase.target, "table entry by %d index", i)
	}
}
