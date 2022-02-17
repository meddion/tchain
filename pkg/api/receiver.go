package api

import (
	"errors"
	"log"

	"github.com/meddion/pkg/crypto"
)

var _ Receiver = &ReceiverRPC{}

var (
	InvalidTransactionError = errors.New("transactions is invalid")
)

type txPool struct {
	txs map[crypto.HashValue]Transaction
}

type ReceiverRPC struct {
	senderPool SenderPool
	txPool     MemPool
	logger     *log.Logger
}

func NewReceiverRPC(memPool MemPool, senderPool SenderPool, logger *log.Logger) Receiver {
	return &ReceiverRPC{senderPool: senderPool, txPool: memPool, logger: logger}
}

func (r *ReceiverRPC) HandleTransaction(req TransactionReq, resp *TransactionResp) error {
	if _, exists := r.txPool.Get(req.Hash); exists {
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

func (r *ReceiverRPC) HandleIsAlive(_ Empty, _ *Empty) error {
	return nil
}

func (r *ReceiverRPC) HandleBlock(req BlockReq, resp *OpStatus) error {
	return nil
}
