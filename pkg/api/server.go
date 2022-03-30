package api

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type Server struct {
	*http.Server
}

func NewServer(rcv Receiver) (*Server, error) {
	s := Server{}
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(rcv); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle(_rpcPath, rpcServer)
	s.Server = &http.Server{
		Handler: mux,
	}

	return &s, nil
}

func (s *Server) Start(addr, port string) error {
	l, err := net.Listen("tcp", addr+":"+port)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *Server) Close(rootCtx context.Context) error {
	ctx, cancel := context.WithTimeout(rootCtx, time.Second*15)
	defer cancel()

	return s.Shutdown(ctx)
}
