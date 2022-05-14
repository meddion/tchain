package core

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/meddion/pkg/crypto"
)

var (
	ErrUnsupportedVer    = errors.New("unsupported version")
	ErrInvalidTimestamp  = errors.New("invalid timestamp")
	ErrInvalidMerkleRoot = errors.New("invalid merkle root")
	ErrInvalidDifficulty = errors.New("invalid difficulty")
	ErrInvalidNonce      = errors.New("invalid nonce")
)

// TODO: impl
func (b Block) Verify() error {
	if !verifyVersion(b.Version) {
		return ErrUnsupportedVer
	}

	if !verifyDifficulty(b.Header.Difficulty) {
		return ErrInvalidDifficulty
	}

	if err := b.Difficulty.VerifyNonce(b.Header); err != nil {
		return ErrInvalidNonce
	}

	for _, tx := range b.Body {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	if hash, err := crypto.GenMerkleRoot(b.Body); err != nil {
		return fmt.Errorf("on generating a merkle root: %w", err)
	} else if b.MerkleRoot != hash {
		return ErrInvalidMerkleRoot
	}

	return nil
}

func verifyDifficulty(diff Difficulty) bool {
	return diff > 14
}

func verifyVersion(ver uint8) bool {
	switch ver {
	case 1:
		return true
	}

	return false
}

var (
	ErrEmptyTxData      = errors.New("empty transaction data")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidChecksum  = errors.New("invalid checksum")
)

func (tx Transaction) Verify() error {
	if len(tx.Data) == 0 {
		return ErrEmptyTxData
	}

	if tx.Sig == nil {
		return ErrInvalidSignature
	}

	hash, err := crypto.Hash256(tx.Data[:])
	if err != nil || !bytes.Equal(hash[:], tx.Hash[:]) {
		return ErrInvalidChecksum
	}

	if !tx.Sig.Verify(hash[:]) {
		return ErrInvalidSignature
	}

	return verifyTransactionData(tx.Data)
}

// TODO: impl
func verifyTransactionData(txData TxData) error {
	return nil
}
