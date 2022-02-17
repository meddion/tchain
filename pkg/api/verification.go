package api

import (
	"bytes"
	"log"

	"github.com/meddion/pkg/crypto"
)

func VerifyTransaction(tx Transaction) bool {
	if tx.R == nil || tx.S == nil || len(tx.Data) == 0 || !tx.PublicKey.IsValid() {
		return false
	}

	hash, err := crypto.Hash(tx.Data)
	// TODO: change the behaviour
	if err != nil || bytes.Compare(hash[:], tx.Hash[:]) != 0 {
		return false
	}

	log.Println(tx)

	return crypto.Verify(tx.PublicKey, hash[:], tx.R, tx.S)
}
