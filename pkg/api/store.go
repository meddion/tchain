package api

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/meddion/pkg/crypto"
)

const (
	_dbPath   = "./blocks.db"
	_dbBucket = "blocks"
)

var (
	ErrBucketNotFound = errors.New("bucket not found")
	ErrMissingBlock   = errors.New("block is missing")
)

type BlockRepo struct {
	db *bolt.DB
}

func NewBlockRepo(db *bolt.DB) *BlockRepo {
	return &BlockRepo{
		db: db,
	}
}

func (b *BlockRepo) Get(blockID crypto.HashValue) (Block, error) {
	var blockBytes []byte
	if err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_dbBucket))
		if b == nil {
			return ErrBucketNotFound
		}

		if blockBytes = b.Get(blockID[:]); len(blockBytes) == 0 {
			return ErrMissingBlock
		}

		return nil
	}); err != nil {
		return Block{}, err
	}

	block, err := DecodeBlock(blockBytes)
	if err != nil {
		return Block{}, err
	}

	return block, nil
}

func (b *BlockRepo) Store(hashKey crypto.HashValue, block Block) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_dbBucket))

		// Block is already stored
		if val := b.Get(hashKey[:]); val != nil {
			return nil
		}

		blockBytes, err := block.Bytes()
		if err != nil {
			return err
		}

		return b.Put(hashKey[:], blockBytes)
	})
}
