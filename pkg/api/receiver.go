package api

import (
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/meddion/pkg/crypto"
)

var _ Receiver = &ReceiverRPC{}

var (
	InvalidTransactionError = errors.New("transactions is invalid")
	InvalidBlockError       = errors.New("block is invalid")
)

const (
	dbPath   = "./blocks.db"
	dbBucket = "blocks"
)

type ReceiverRPC struct {
	senderPool SenderPool
	txPool     map[crypto.HashValue]Transaction
	db         *bolt.DB
	logger     *log.Logger
}

func NewReceiverRPC(senderPool SenderPool, logger *log.Logger) Receiver {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalln(err)
	}

	return &ReceiverRPC{senderPool: senderPool, db: db, logger: logger}
}

func (r *ReceiverRPC) HandleTransaction(req TransactionReq, resp *TransactionResp) error {
	if _, exists := r.txPool[req.Hash]; exists {
		return nil
	}

	if !VerifyTransaction(req.Transaction) {
		return InvalidTransactionError
	}

	// TODO: handle errors
	var errors []error
	for _, s := range r.senderPool.Senders() {
		if _, err := s.SendTransaction(req); err != nil {
			errors = append(errors, err)
		}
	}
	// TODO: change it
	r.logger.Println(errors)

	return nil
}

func (r *ReceiverRPC) HandleBlock(req BlockReq, resp *OpStatus) error {
	if !VerifyBlock(req.Block) {
		return InvalidBlockError
	}

	hbytes, err := req.Block.Header.Bytes()
	if err != nil {
		return err
	}

	hashKey, err := crypto.Hash256(hbytes)
	if err != nil {
		return err
	}

	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))

		blockBytes, err := req.Block.Bytes()
		if err != nil {
			return err
		}

		return b.Put(hashKey[:], blockBytes)
	})
}

func (r *ReceiverRPC) HandleIsAlive(_ Empty, _ *Empty) error {
	return nil
}
