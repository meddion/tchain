package api

import (
	"errors"
	"log"

	"github.com/meddion/pkg/crypto"
)

var _ Receiver = &ReceiverRPC{}

var (
	ErrInvalidTransaction = errors.New("transactions is invalid")
	ErrInvalidBlock       = errors.New("block is invalid")
)

type ReceiverRPC struct {
	txPool     map[crypto.HashValue]Transaction
	senderPool PeerPool
	db         *BlockRepo

	logger *log.Logger
}

func NewReceiverRPC(senderPool PeerPool, db *BlockRepo, logger *log.Logger) Receiver {
	return &ReceiverRPC{
		txPool:     make(map[crypto.HashValue]Transaction),
		senderPool: senderPool,
		db:         db,
		logger:     logger,
	}
}

func (r *ReceiverRPC) HandleTransaction(req TransactionReq, resp *TransactionResp) error {
	if _, exists := r.txPool[req.Hash]; exists {
		return nil
	}

	if err := VerifyTransaction(req.Transaction); err != nil {
		return err
	}

	r.txPool[req.Hash] = req.Transaction

	// TODO: handle errors
	var errors []error
	for _, s := range r.senderPool.Peers() {
		if _, err := s.SendTransaction(req); err != nil {
			errors = append(errors, err)
		}
	}
	// TODO: change it
	r.logger.Println(errors)

	return nil
}

func (r *ReceiverRPC) HandleBlock(req BlockReq, resp *Empty) error {
	prevBlock, err := r.db.Get(req.Block.PrevBlockHash)
	if err != nil {
		return err
	}

	if err := VerifyBlock(req.Block, prevBlock); err != nil {
		return err
	}

	hbytes, err := req.Block.Header.Bytes()
	if err != nil {
		return err
	}

	hashKey, err := crypto.Hash256(hbytes)
	if err != nil {
		return err
	}

	return r.db.Store(hashKey, req.Block)
}

func (r *ReceiverRPC) HandleIsAlive(_ Empty, _ *Empty) error {
	return nil
}
