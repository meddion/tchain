package core

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/meddion/pkg/crypto"
)

type Blockchain struct {
	db     *BlockRepo
	logger *log.Logger

	index blockIndex

	mtx      sync.RWMutex
	lastNode *blockNode

	// currentView blockView
	// lock    sync.RWMutex
	// orphans map[crypto.HashValueValue]*blockNode
}

func NewBlockchain(db *BlockRepo, logger *log.Logger) (*Blockchain, error) {
	b := &Blockchain{
		db:     db,
		logger: logger,
		index:  newBlockIndex(),
	}

	if err := b.setGenesisBlock(getGenesisPair()); err != nil {
		return nil, fmt.Errorf("on setting a genesis block: %w", err)
	}

	return b, nil
}

func (b *Blockchain) setGenesisBlock(hash crypto.HashValue, block Block) error {
	node, err := newBlockNode(nil, block.Header)
	if err != nil {
		return err
	}

	b.index.AddNode(node)

	if err := b.db.Store(hash[:], block); err != nil {
		return err
	}

	b.setLastBlock(node)
	if err := b.db.Store(_lastCommitedBlockNodeKey, node); err != nil {
		return err
	}

	return nil
}

// func (b *Blockchain) getOrphan(hashKey crypto.HashValueValue) (n *blockNode, ok bool) {
// 	b.lock.RLock()
// 	defer b.lock.RUnlock()

// 	n, ok = b.orphans[hashKey]
// 	return
// }

// func (b *Blockchain) setOrphan(hashKey crypto.HashValueValue, node *blockNode) {
// 	b.lock.Lock()
// 	defer b.lock.Unlock()

// 	b.orphans[hashKey] = node
// }

func (b *Blockchain) ProcessBlock(block Block) error {
	hashKey, err := block.Header.Checksum()
	if err != nil {
		return err
	}

	if n := b.index.GetNode(hashKey); n != nil {
		return errors.New("dublicate block")
	}

	// if _, err := b.db.Get(hashKey); err == nil {
	// }

	// if _, exists := b.getOrphan(hashKey); exists {
	// 	return errors.New("dublicate orphan block")
	// }

	// prevBlock, err := b.db.Get(block.PrevBlockHash)
	// if err != nil {
	// 	// b.setOrphan(hashKey, &blockNode{blk: &block})
	// 	return err
	// } else if block.Timestamp < prevBlock.Timestamp {
	// 	return ErrInvalidTimestamp
	// }

	parentNode := b.index.GetNode(block.PrevBlockHash)
	// TODO: add to orphans
	if parentNode == nil {
		return ErrMissingParentNode
	}

	if block.Timestamp < parentNode.Timestamp {
		return ErrInvalidTimestamp
	}

	if err := block.Verify(); err != nil {
		return err
	}

	node, err := newBlockNode(parentNode, block.Header)
	if err != nil {
		return err
	}

	b.index.AddNode(node)

	if err := b.db.Store(hashKey[:], block); err != nil {
		return err
	}

	return b.connectNodeToChain(node)
}

func (b *Blockchain) setLastBlock(node *blockNode) {
	b.mtx.Lock()
	b.lastNode = node
	b.mtx.Unlock()
}

func (b *Blockchain) connectNodeToChain(node *blockNode) error {
	// Adding to the tip
	if b.lastNode.Hash == node.Prev.Hash {
		b.setLastBlock(node)
		if err := b.db.Store(_lastCommitedBlockNodeKey, node); err != nil {
			return err
		}

		return nil
	}

	// Adding to an existing chain
	// If the amound of work on this chain is larger -- make it the man chain

	return nil
}
