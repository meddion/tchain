package tchain

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/rpc"
)

type SrvRPC struct {
}

type Message string

func (s *SrvRPC) Ping(req Message, resp *Message) error {
	if string(req) != "ping" {
		return errors.New("expected ping")
	}
	*resp = Message("pong")

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

	s.Server = &http.Server{
		Handler: rpcServer,
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
