package core

import (
	"log"

	"github.com/meddion/pkg/crypto"
)

var _ Receiver = &ReceiverRPC{}

type ReceiverRPC struct {
	blkchain *Blockchain
	peerPool PeerPool
	txPool   map[crypto.HashValue]Transaction
	logger   *log.Logger
}

func NewReceiverRPC(blkchain *Blockchain, senderPool PeerPool, logger *log.Logger) Receiver {
	return &ReceiverRPC{
		blkchain: blkchain,
		txPool:   make(map[crypto.HashValue]Transaction),
		peerPool: senderPool,
		logger:   logger,
	}
}

type PeerPool interface {
	NumberOfPeers() int
	SendToPeers(func(Peer) error) <-chan error
	Add(Peer)
	Peers() []Peer
	Close()
}

func (r *ReceiverRPC) propagateToPeers(f func(s Peer) error) {
	for err := range r.peerPool.SendToPeers(f) {
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

	r.propagateToPeers(func(p Peer) error {
		_, err := p.SendTransaction(req)
		return err
	})

	return nil
}

func (r *ReceiverRPC) HandleBlock(req BlockReq, resp *Empty) error {
	if err := r.blkchain.ProcessBlock(req.Block); err != nil {
		return err
	}

	r.propagateToPeers(func(p Peer) error {
		return p.SendBlock(req)
	})

	return nil
}

func (r *ReceiverRPC) HandleIsAlive(_ Empty, _ *Empty) error {
	return nil
}

func (r *ReceiverRPC) HandlePeersDiscovery(_ Empty, knownPeers *PeersDiscoveryResp) error {
	peers := r.peerPool.Peers()
	addrs := make([]Addr, len(peers))
	for i, p := range peers {
		addrs[i] = p.Addr()
	}
	knownPeers.addrs = addrs

	return nil
}
