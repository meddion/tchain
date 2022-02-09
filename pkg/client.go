package tchain

import (
	"net/rpc"
)

type Client struct {
	*rpc.Client
}

func NewClient(addr, port string) (Client, error) {
	c, err := rpc.DialHTTPPath("tcp", addr+":"+port, RPCpath)
	if err != nil {
		return Client{}, err
	}
	client := Client{c}

	return client, nil
}

func (c Client) Ping(msg TxRequest) (TxRequest, error) {
	var reply TxRequest
	err := c.Call("SrvRPC.Transaction", &msg, &reply)
	if err != nil {
		return TxRequest(""), err
	}

	return reply, nil
}
