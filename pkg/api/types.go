package api

import (
	"math/big"

	"github.com/meddion/pkg/crypto"
)

type (
	MemPool interface {
		Get(crypto.HashValue) (Transaction, bool)
		Store(Transaction)
	}

	SenderPool interface {
		Senders() []Sender
	}
)

type (
	Header struct {
		version    int8
		Timestamp  int64
		Hash       [32]byte
		MerkleRoot [32]byte
		nonce      int32
	}

	Body [1_000_000]byte

	Block struct {
		Header
		Body
	}
)

type (
	Transaction struct {
		PublicKey crypto.PublicKey
		R, S      *big.Int
		Hash      crypto.HashValue
		Data      []byte
	}
)
