package core

import (
	"bytes"
	"encoding/gob"
	"math"

	"github.com/meddion/pkg/crypto"
)

const (
	_rpcPath = "/_tchain_rpc_"
)

type Sender interface {
	SendTransaction(TransactionReq) (TransactionResp, error)
	SendIsAlive() error
	SendBlock(BlockReq) error
	SendPeersDiscovery() (PeersDiscoveryResp, error)
}

type Receiver interface {
	HandleTransaction(TransactionReq, *TransactionResp) error
	HandleIsAlive(Empty, *Empty) error
	HandleBlock(BlockReq, *Empty) error
	HandlePeersDiscovery(Empty, *PeersDiscoveryResp) error
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

	PeersDiscoveryResp struct {
		addrs []Addr
	}
)
type Addr struct {
	IP, Port string
}

func (a Addr) String() string {
	return a.IP + ":" + a.Port
}

const (
	NonceMaxValue    = math.MaxUint32
	BodyElementLimit = 64
	TxBodySizeLimit  = 1024
)

type Signature interface {
	Verify([]byte) bool
}

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
		Difficulty    Difficulty
		Nonce         Nonce
	}

	Body   = []Transaction
	TxData []byte

	Transaction struct {
		Sig  Signature
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
