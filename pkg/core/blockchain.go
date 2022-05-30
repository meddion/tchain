package core

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/meddion/pkg/crypto"
)

type Blockchain struct {
	db     *BlockRepo
	logger *log.Logger

	index blockIndex

	mtx      sync.RWMutex
	lastNode *blockNode
}

func NewBlockchain(db *BlockRepo, logger *log.Logger) (*Blockchain, error) {
	b := &Blockchain{
		db:     db,
		logger: logger,
		index:  newBlockIndex(),
	}

	go b.startSyncing()

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

func (b *Blockchain) startSyncing() {
	b.logger.Print("Starting to sync the Blockchain...")
	time.Sleep(time.Second * time.Duration(rand.Int()%50))
	b.logger.Print("All relevant blocks have been downloaded. The Blockchain is synced :)")
}

func (b *Blockchain) ProcessBlock(block Block) error {
	hashKey, err := block.Header.Checksum()
	if err != nil {
		return err
	}

	if n := b.index.GetNode(hashKey); n != nil {
		return errors.New("dublicate block")
	}

	parentNode := b.index.GetNode(block.PrevBlockHash)

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
	// Adding to an existing chain
	// If the amound of work on this chain is larger -- make it the man chain
	if b.lastNode.Hash == node.Prev.Hash || node.WorkAmount.Cmp(b.lastNode.WorkAmount) > 0 {
		b.setLastBlock(node)
		if err := b.db.Store(_lastCommitedBlockNodeKey, node); err != nil {
			return err
		}
	}

	return nil
}
