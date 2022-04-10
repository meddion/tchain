package api

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/meddion/pkg/crypto"
)

var (
	GenesisBlock     Block
	GenesisBlockHash crypto.HashValue
)

func init() {
	msg := "on creating a genesis block: %v"

	ghash, err := crypto.Hash256([]byte("genesis"))
	if err != nil {
		log.Fatalf(msg, err)
	}

	mroot, err := crypto.GenMerkleRoot([]Transaction{})
	if err != nil {
		log.Fatalf(msg, err)
	}

	GenesisBlock = Block{
		Header: Header{
			Version:       1,
			Timestamp:     time.Date(2021, time.February, 24, 0, 0, 0, 0, time.UTC).Unix(),
			PrevBlockHash: ghash,
			MerkleRoot:    mroot,
			Nonce:         0,
		},
		Body: Body{},
	}

	b, err := GenesisBlock.Bytes()
	if err != nil {
		log.Fatalf(msg, err)
	}

	GenesisBlockHash, err = crypto.Hash256(b)
	if err != nil {
		log.Fatalf(msg, err)
	}
}

const (
	_rpcPath = "/_tchain_rpc_"
)

type Sender interface {
	SendTransaction(TransactionReq) (TransactionResp, error)
	SendIsAlive() error
	SendBlock(BlockReq) error
}

type Receiver interface {
	HandleTransaction(TransactionReq, *TransactionResp) error
	HandleIsAlive(Empty, *Empty) error
	HandleBlock(BlockReq, *Empty) error
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

type PeerPool interface {
	Peers() []Sender
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
		Version       uint8
		Timestamp     int64
		PrevBlockHash crypto.HashValue
		MerkleRoot    crypto.HashValue
		Nonce         uint32
	}

	Body   = []Transaction
	TxData []byte

	Transaction struct {
		Sig  crypto.Signature
		Hash crypto.HashValue
		Data TxData
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

func (b *Block) FromBytes(data []byte) error {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&b); err != nil {
		return err
	}

	return nil
}

func (t Transaction) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(t); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
