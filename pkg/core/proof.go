package core

import (
	"errors"
	"math/big"

	"github.com/meddion/pkg/crypto"
)

const _blockDifficalty = 21

var (
	_bigOne    = big.NewInt(1)
	_oneLsh256 = new(big.Int).Lsh(_bigOne, 256)
)

type (
	Nonce      uint32
	Difficulty uint32
)

func (d Difficulty) DifficultyBits() *big.Int {
	i := big.NewInt(1)
	i.Lsh(i, uint(255-d))

	return i
}

func (d Difficulty) WorkAmount() *big.Int {
	denominator := new(big.Int).Add(d.DifficultyBits(), _bigOne)

	return new(big.Int).Div(_oneLsh256, denominator)
}

func (d Difficulty) GenNonce(header Header) (Nonce, error) {
	var (
		tempInt big.Int
		target  = d.DifficultyBits()
	)

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

func (d Difficulty) VerifyNonce(header Header) error {
	hb, err := header.Bytes()
	if err != nil {
		return err
	}

	hash, err := crypto.Hash256(hb)
	if err != nil {
		return err
	}

	var tempInt big.Int
	tempInt.SetBytes(hash[:])
	if tempInt.Cmp(d.DifficultyBits()) == -1 {
		return nil
	}

	return errors.New("verifying nonce againts target")
}
