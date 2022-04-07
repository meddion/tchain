package api

import (
	"net/rpc"
)

var _ Sender = &SenderRPC{}

type SenderRPC struct {
	client *rpc.Client
}

func NewSender(addr, port string) (Sender, error) {
	c, err := rpc.DialHTTPPath("tcp", addr+":"+port, _rpcPath)
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
	// TODO: get it through reflect
	err := s.client.Call("ReceiverRPC.HandleIsAlive", &Empty{}, &Empty{})

	if err != nil {
		return err
	}

	return nil
}

func (s SenderRPC) SendBlock(blockReq BlockReq) error {
	var opStat OpStatus
	err := s.client.Call("ReceiverRPC.HandleBlock", blockReq, opStat)
	if err != nil {
		return err
	}

	return nil
}
