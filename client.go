package tchain

import (
	"net/rpc"
)

type Client struct {
	*rpc.Client
}

func NewClient(addr, port string) (Client, error) {
	c, err := rpc.DialHTTP("tcp", addr+":"+port)
	if err != nil {
		return Client{}, err
	}
	client := Client{c}

	return client, nil
}

func (c Client) Ping(msg Message) (Message, error) {
	var reply Message
	err := c.Call("SrvRPC.Ping", &msg, &reply)
	if err != nil {
		return Message(""), err
	}

	return reply, nil
}
