package core

import (
	"bytes"
	"encoding/gob"
	"math/big"
	"sync"

	"github.com/meddion/pkg/crypto"
)

// Fields are exported to satisfy gob encoder/decorer
type blockNode struct {
	Prev       *blockNode
	WorkAmount *big.Int
	Height     int
	Hash       crypto.HashValue

	// Some fields from Header
	Version    uint8
	Timestamp  int64
	MerkleRoot crypto.HashValue
	Nonce      Nonce
}

func newBlockNode(prev *blockNode, header Header) (*blockNode, error) {
	hash, err := header.Checksum()
	if err != nil {
		return nil, err
	}

	node := &blockNode{
		Hash:       hash,
		WorkAmount: header.Difficulty.WorkAmount(),
		Version:    header.Version,
		Timestamp:  header.Timestamp,
		MerkleRoot: header.MerkleRoot,
		Nonce:      header.Nonce,
	}
	if prev != nil {
		node.Prev = prev
		node.Height = prev.Height + 1
		node.WorkAmount.Add(node.WorkAmount, node.Prev.WorkAmount)
	}

	return node, nil
}

func (node *blockNode) Ancestor(height int) *blockNode {
	if height < 0 || height > node.Height {
		return nil
	}

	n := node
	for ; n != nil && n.Height != height; n = n.Prev {
	}

	return n
}

func (b blockNode) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type blockIndex struct {
	mut   sync.RWMutex
	index map[crypto.HashValue]*blockNode
}

func newBlockIndex() blockIndex {
	return blockIndex{
		index: make(map[crypto.HashValue]*blockNode),
	}
}

func (b *blockIndex) IsNodePresent(hash crypto.HashValue) bool {
	b.mut.RLock()
	defer b.mut.RUnlock()

	_, exists := b.index[hash]
	return exists
}

func (b *blockIndex) GetNode(hash crypto.HashValue) *blockNode {
	b.mut.RLock()
	defer b.mut.RUnlock()

	return b.index[hash]
}

func (b *blockIndex) AddNode(node *blockNode) {
	b.mut.Lock()
	b.index[node.Hash] = node
	b.mut.Unlock()
}
