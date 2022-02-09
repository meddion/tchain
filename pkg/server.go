package tchain

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
)

const RPCpath = "/_tchain_rpc"

type SrvRPC struct {
}

type TxRequest struct {
	ttl  int
	data []byte
}

func (s *SrvRPC) Transaction(req TxRequest, resp *TxRequest) error {
	if req.ttl--; req.ttl <= 0 {
		return nil
	}

	// Call other clients
	for _, c := range clients {
		c.Transaction(req)
	}

	return nil
}

type Service struct {
	*http.Server
}

func NewService() (*Service, error) {
	s := Service{}
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(&SrvRPC{}); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle(RPCpath, rpcServer)
	s.Server = &http.Server{
		Handler: mux,
	}

	return &s, nil
}

func (s *Service) Start(addr, port string) error {
	l, err := net.Listen("tcp", addr+":"+port)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *Service) Stop() error {
	return s.Shutdown(context.TODO())
}
