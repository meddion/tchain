package core

import (
	"errors"
	"net/rpc"
	"time"
)

const _isAliveWaitDuration = time.Second * 5

var ErrIsAliveTimeout = errors.New("timeout for peer")

var _ Sender = &SenderRPC{}

type SenderRPC struct {
	client *rpc.Client
}

func NewSender(addr Addr) (Sender, error) {
	c, err := rpc.DialHTTPPath("tcp", addr.String(), _rpcPath)
	if err != nil {
		return SenderRPC{}, err
	}

	return SenderRPC{client: c}, nil
}

func (s SenderRPC) SendTransaction(req TransactionReq) (TransactionResp, error) {
	var resp TransactionResp

	err := s.client.Call("ReceiverRPC.HandleTransaction", &req, &resp)

	if err != nil {
		return TransactionResp{}, err
	}

	return resp, nil
}

func (s SenderRPC) SendIsAlive() error {
	call := s.client.Go("ReceiverRPC.HandleIsAlive", &Empty{}, &Empty{}, make(chan *rpc.Call, 1))

	select {
	case c := <-call.Done:
		return c.Error
	case <-time.After(_isAliveWaitDuration):
	}

	return ErrIsAliveTimeout
}

func (s SenderRPC) SendBlock(blockReq BlockReq) error {
	err := s.client.Call("ReceiverRPC.HandleBlock", blockReq, &Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (s SenderRPC) SendPeersDiscovery() (PeersDiscoveryResp, error) {
	var knownPeers PeersDiscoveryResp
	err := s.client.Call("ReceiverRPC.HandleBlock", Empty{}, &knownPeers)
	if err != nil {
		return PeersDiscoveryResp{}, err
	}

	return knownPeers, nil
}
