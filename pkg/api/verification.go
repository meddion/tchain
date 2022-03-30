package api

import (
	"bytes"

	"github.com/meddion/pkg/crypto"
)

// TODO: impl
func VerifyBlock(block, prevBlock Block) bool {
	if !verifyVersion(block.version) {
		return false
	}

	if block.Timestamp < prevBlock.Timestamp {
		return false
	}

	for _, tx := range block.Body {
		if !VerifyTransaction(tx) {
			// TODO: log

			return false
		}
	}

	byteArrays, err := block.ByteArrays()
	if err != nil {
		return false
	}

	hash, err := crypto.GenMerkleRoot(byteArrays)
	if err != nil {
		return false
	}

	return block.MerkleRoot == hash
}

func verifyVersion(ver uint8) bool {
	switch ver {
	case 1:
		return true
	}

	return false
}

func VerifyTransaction(tx Transaction) bool {
	if len(tx.Data) == 0 || tx.Sig == nil {
		return false
	}

	hash, err := crypto.Hash256(tx.Data[:])
	if err != nil || !bytes.Equal(hash[:], tx.Hash[:]) {
		return false
	}

	if !tx.Sig.Verify(hash[:]) {
		return false
	}

	return VerifyTransactionData(tx.Data)
}

// TODO: impl
func VerifyTransactionData(txData TxData) bool {
	return true
}
