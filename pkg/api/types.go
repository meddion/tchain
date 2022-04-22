package api

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"time"

	"github.com/meddion/pkg/crypto"
)

const (
	_rpcPath = "/_tchain_rpc_"
)

var (
	_genesisBlock     Block
	_genesisBlockHash crypto.HashValue
)

func getGenesisPair() (crypto.HashValue, Block) {
	if _genesisBlock.Nonce != 0 {
		return _genesisBlockHash, _genesisBlock
	}

	var err error
	if _genesisBlock, err = genGenesisBlock(); err != nil {
		panic(fmt.Sprintf("on creating a genesis block: %v", err))
	}

	if _genesisBlockHash, err = _genesisBlock.Header.Checksum(); err != nil {
		panic(fmt.Sprintf("on hashing a genesis block: %v", err))
	}

	return _genesisBlockHash, _genesisBlock
}

func genGenesisBlock() (Block, error) {
	ghash, err := crypto.Hash256([]byte("genesis"))
	if err != nil {
		return Block{}, err
	}

	mroot, err := crypto.GenMerkleRoot([]Transaction{})
	if err != nil {
		return Block{}, err
	}

	header := Header{
		Version:       1,
		Timestamp:     time.Date(2021, time.February, 24, 0, 0, 0, 0, time.UTC).Unix(),
		PrevBlockHash: ghash,
		MerkleRoot:    mroot,
		Nonce:         2068160, // difficalty = 21
	}
	// nonce, err := genPowNonce(header, getPowTarget())
	// if err != nil {
	// 	return Block{}, err

	// }
	// header.Nonce = nonce

	return Block{
		Header: header,
		Body:   Body{},
	}, nil
}

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
	NonceMaxValue    = math.MaxUint32
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
		Nonce         Nonce
	}

	Nonce  uint32
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

func (h Header) Checksum() (crypto.HashValue, error) {
	b, err := h.Bytes()
	if err != nil {
		return crypto.ZeroHashValue, err
	}

	headerHash, err := crypto.Hash256(b)
	if err != nil {
		return crypto.ZeroHashValue, err
	}

	return headerHash, nil
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
