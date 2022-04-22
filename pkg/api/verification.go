package api

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/meddion/pkg/crypto"
)

var (
	ErrUnsupportedVer    = errors.New("unsupported version")
	ErrInvalidTimestamp  = errors.New("invalid timestamp")
	ErrInvalidMerkleRoot = errors.New("invalid merkle root")
	ErrInvalidNonce      = errors.New("invalid nonce")
)

// TODO: impl
func VerifyBlock(block, prevBlock Block) error {
	if !verifyVersion(block.Version) {
		return ErrUnsupportedVer
	}

	if block.Timestamp < prevBlock.Timestamp {
		return ErrInvalidTimestamp
	}

	if err := verifyHeaderNonce(block.Header, getPowTarget()); err != nil {
		return err
	}

	for _, tx := range block.Body {
		if err := VerifyTransaction(tx); err != nil {
			return err
		}
	}

	if hash, err := crypto.GenMerkleRoot(block.Body); err != nil {
		return fmt.Errorf("on generating a merkle root: %w", err)
	} else if block.MerkleRoot != hash {
		return ErrInvalidMerkleRoot
	}

	return nil
}

func verifyHeaderNonce(header Header, target powTarget) error {
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
	if tempInt.Cmp(target) == -1 {
		return nil
	}

	return ErrInvalidNonce
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

func VerifyTransaction(tx Transaction) error {
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

	return VerifyTransactionData(tx.Data)
}

// TODO: impl
func VerifyTransactionData(txData TxData) error {
	return nil
}
