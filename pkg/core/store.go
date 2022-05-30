package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/meddion/pkg/crypto"
)

const (
	_dbPath   = "./blocks.db"
	_dbBucket = "blocks"
)

var _lastCommitedBlockNodeKey = []byte("lastCommited")

var (
	ErrBucketNotFound    = errors.New("bucket not found")
	ErrMissingBlock      = errors.New("block is missing")
	ErrMissingParentNode = errors.New("parent node is missing")
)

type BlockRepo struct {
	db *bolt.DB
}

func NewBlockRepo(dbFile string) (*BlockRepo, error) {
	if dbFile == "" {
		dbFile = _dbPath
	}

	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("on opening a bolt conn: %w", err)
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(_dbBucket))
		if err != nil && err != bolt.ErrBucketExists {
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

	var block Block
	if err := block.FromBytes(blockBytes); err != nil {
		return Block{}, err
	}

	return block, nil
}

type bytesConverter interface {
	Bytes() ([]byte, error)
}

func (b *BlockRepo) Store(hashKey []byte, block bytesConverter) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_dbBucket))
		if b == nil {
			return ErrBucketNotFound
		}

		// Block is already stored
		if val := bucket.Get(hashKey); val != nil {
			return nil
		}

		blockBytes, err := block.Bytes()
		if err != nil {
			return err
		}

		return bucket.Put(hashKey, blockBytes)
	})
}

func (b *BlockRepo) Close() error {
	return b.db.Close()
}
