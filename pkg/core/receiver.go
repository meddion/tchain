package core

import (
	"log"

	"github.com/meddion/pkg/crypto"
)

var _ Receiver = &ReceiverRPC{}

type ReceiverRPC struct {
	blkchain   *Blockchain
	senderPool PeerPool
	txPool     map[crypto.HashValue]Transaction
	logger     *log.Logger
}

func NewReceiverRPC(blkchain *Blockchain, senderPool PeerPool, logger *log.Logger) Receiver {
	return &ReceiverRPC{
		blkchain:   blkchain,
		txPool:     make(map[crypto.HashValue]Transaction),
		senderPool: senderPool,
		logger:     logger,
	}
}

type PeerPool interface {
	NumberOfPeers() int
	SendToPeers(func(Sender) error) <-chan error
	Peers() []Peer
}

func (r *ReceiverRPC) propagateToPeers(s func(s Sender) error) {
	for err := range r.senderPool.SendToPeers(s) {
		r.logger.Println(err)
	}
}

func (r *ReceiverRPC) HandleTransaction(req TransactionReq, resp *TransactionResp) error {
	if _, exists := r.txPool[req.Hash]; exists {
		return nil
	}

	if err := req.Verify(); err != nil {
		return err
	}

	r.txPool[req.Hash] = req.Transaction

	r.propagateToPeers(func(s Sender) error {
		_, err := s.SendTransaction(req)
		return err
	})

	return nil
}

func (r *ReceiverRPC) HandleBlock(req BlockReq, resp *Empty) error {
	if err := r.blkchain.ProcessBlock(req.Block); err != nil {
		return err
	}

	r.propagateToPeers(func(s Sender) error {
		return s.SendBlock(req)
	})

	return nil
}

func (r *ReceiverRPC) HandleIsAlive(_ Empty, _ *Empty) error {
	return nil
}
