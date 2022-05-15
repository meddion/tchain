package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPowGenesis(t *testing.T) {
	pow := Difficulty(15)
	h := Header{}
	nonce, err := pow.GenNonce(h)
	assert.NoError(t, err)

	h.Nonce = nonce
	assert.NoError(t, pow.VerifyNonce(h))
}
