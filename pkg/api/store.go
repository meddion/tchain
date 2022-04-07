package api

import (
	"errors"
	"fmt"

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

type BlockPair struct {
	hash crypto.HashValue
	Block
}

type BlockRepo struct {
	db           *bolt.DB
	lastCommited []BlockPair
}

func NewBlockRepo(db *bolt.DB) (*BlockRepo, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket([]byte(_dbBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &BlockRepo{
		db: db,
	}, nil
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
		bucket := tx.Bucket([]byte(_dbBucket))
		if b == nil {
			return ErrBucketNotFound
		}

		// Block is already stored
		if val := bucket.Get(hashKey[:]); val != nil {
			return nil
		}

		blockBytes, err := block.Bytes()
		if err != nil {
			return err
		}

		// TODO: change caching
		b.cacheBlock(hashKey, block)

		return bucket.Put(hashKey[:], blockBytes)
	})
}

func (b *BlockRepo) cacheBlock(hash crypto.HashValue, blk Block) {
	b.lastCommited = append(b.lastCommited, BlockPair{hash, blk})
}

func (b *BlockRepo) LastCommited() BlockPair {
	return b.lastCommited[len(b.lastCommited)-1]
}

func (b *BlockRepo) Close() error {
	return b.db.Close()
}
