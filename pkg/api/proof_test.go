package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPowGenesis(t *testing.T) {
	target := newPowTarget(19)
	h := Header{}
	nonce, err := genPowNonce(h, target)
	assert.NoError(t, err)

	h.Nonce = nonce
	assert.NoError(t, verifyHeaderNonce(h, target))
}
