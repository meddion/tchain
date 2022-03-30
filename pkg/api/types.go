package api

import (
	"bytes"
	"encoding/gob"

	"github.com/meddion/pkg/crypto"
)

const (
	_rpcPath = "/_tchain_rpc_"
)

type Sender interface {
	SendTransaction(TransactionReq) (TransactionResp, error)
	SendIsAlive() error
	SendBlock() error
}

type Receiver interface {
	HandleTransaction(TransactionReq, *TransactionResp) error
	HandleIsAlive(Empty, *Empty) error
	HandleBlock(BlockReq, *OpStatus) error
}

type (
	Empty struct{}

	OpStatus struct {
		Status bool
		Msg    string
	}

	BlockReq struct {
		Block
	}

	TransactionReq struct {
		Transaction
	}

	TransactionResp struct {
		Status bool
		Msg    string
	}
)

type SenderPool interface {
	Senders() []Sender
}

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
		Sig  crypto.Signature
		Hash crypto.HashValue
		Data TxData
	}
)

// TODO: Rewrite using Go generics

func (h Header) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(h); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t Transaction) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(t); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b Body) ByteArrays() ([][]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	byteArrays := make([][]byte, 0, BodyElementLimit)

	for _, tx := range b {
		txBytes, err := tx.Bytes()
		if err != nil {
			return nil, err
		}

		byteArrays = append(byteArrays, txBytes)
	}

	if err := enc.Encode(b); err != nil {
		return nil, err
	}

	return byteArrays, nil
}

func (b Block) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeBlock(data []byte) (Block, error) {
	var block Block
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&block); err != nil {
		return Block{}, err
	}

	return block, nil
}
