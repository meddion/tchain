package api

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/meddion/pkg/crypto"
)

type (
	SenderPool interface {
		Senders() []Sender
	}
)

const (
	BodyElementLimit = 64
	TxBodySizeLimit  = 1024
)

type (
	Block struct {
		Header
		Body
	}

	Header struct {
		version       uint8
		Timestamp     int64
		PrevBlockHash crypto.HashValue
		MerkleRoot    crypto.HashValue
		Nonce         uint32
	}

	Body   [BodyElementLimit]Transaction
	TxData [TxBodySizeLimit]byte

	Transaction struct {
		PublicKey crypto.PublicKey
		R, S      *big.Int
		Hash      crypto.HashValue
		Data      TxData
	}
)

func (h Header) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(h); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b Block) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decode(data []byte) (Block, error) {
	var block Block
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&block); err != nil {
		return Block{}, err
	}

	return block, nil
}
