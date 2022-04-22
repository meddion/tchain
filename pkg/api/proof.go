package api

import (
	"errors"
	"math/big"

	"github.com/meddion/pkg/crypto"
)

const _blockDifficalty = 21

type powTarget *big.Int

var _powTarget = newPowTarget(_blockDifficalty)

func getPowTarget() powTarget {
	return _powTarget
}

func newPowTarget(difficalty int) powTarget {
	i := big.NewInt(1)
	i.Lsh(i, uint(255-difficalty))

	return powTarget(i)
}

func genPowNonce(header Header, target powTarget) (Nonce, error) {
	var tempInt big.Int

	for i := Nonce(0); i < NonceMaxValue; i++ {
		header.Nonce = i
		hb, err := header.Bytes()
		if err != nil {
			return 0, err
		}
		hash, err := crypto.Hash256(hb)
		if err != nil {
			return 0, err
		}
		tempInt.SetBytes(hash[:])

		if tempInt.Cmp(target) == -1 {
			return i, nil
		}
	}

	return 0, errors.New("pow not found")
}
