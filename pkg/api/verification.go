package api

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/meddion/pkg/crypto"
)

// TODO: impl
func VerifyBlock(db *bolt.DB, block Block) bool {
	if !VerifyVersion(block.version) {
		return false
	}

	var prevBlockBytes []byte
	// Find the previous block
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			panic("bucket doesn't exist")
		}

		prevBlockBytes = b.Get(block.PrevBlockHash[:])
		if prevBlockBytes == nil {
			return errors.New("block doesn't exist")
		}

		return nil
	}); err != nil {
		return false
	}

	prevBlock, err := DecodeBlock(prevBlockBytes)
	if err != nil {
		return false
	}

	if block.Timestamp < prevBlock.Timestamp {
		return false
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

func VerifyVersion(ver uint8) bool {
	switch ver {
	case 1:
		return true
	}

	return false
}

func VerifyTransaction(tx Transaction) bool {
	if tx.R == nil || tx.S == nil || len(tx.Data) == 0 || !tx.PublicKey.IsValid() {
		return false
	}

	hash, err := crypto.Hash256(tx.Data[:])
	// TODO: change the behaviour
	if err != nil || bytes.Compare(hash[:], tx.Hash[:]) != 0 {
		return false
	}

	if !crypto.Verify(tx.PublicKey, hash[:], tx.R, tx.S) {
		return false
	}

	return VerifyTransactionData(tx.Data)
}

// TODO: impl
func VerifyTransactionData(txData TxData) bool {
	return true
}
