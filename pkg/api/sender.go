package api

import (
	"net/rpc"
)

var _ Sender = &SenderRPC{}

type SenderRPC struct {
	*rpc.Client
	pool SenderPool
}

func NewSender(pool SenderPool, addr, port string) (Sender, error) {
	c, err := rpc.DialHTTPPath("tcp", addr+":"+port, RpcPath)
	if err != nil {
		return SenderRPC{}, err
	}

	return SenderRPC{Client: c, pool: pool}, nil
}

func (s SenderRPC) SendTransaction(req TransactionReq) (TransactionResp, error) {
	var resp TransactionResp

	err := s.Call("ReceiverRPC.HandleTransaction", &req, &resp)

	if err != nil {
		return TransactionResp{}, err
	}

	return resp, nil
}

func (s SenderRPC) SendIsAlive() error {
	// TODO: get it through reflect
	err := s.Call("ReceiverRPC.HandleIsAlive", &Empty{}, &Empty{})

	if err != nil {
		return err
	}

	return nil
}

func (s SenderRPC) SendBlock() error {
	// TODO: get it through reflect
	// err := s.Call("ReceiverRPC.HandleBlock", &BlockReq{}, &OpStatus{})

	// if err != nil {
	// 	return err
	// }

	return nil
}
