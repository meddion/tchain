package tchain

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const port = "2022"
const addr = ""

func TestRPC(t *testing.T) {
	s, err := NewService()
	assert.NoError(t, err)
	go func() { assert.NoError(t, s.Start(addr, port)) }()

	<-time.After(time.Second * 2)

	c, err := NewClient(addr, port)
	assert.NoError(t, err)

	resp, err := c.Ping(Message("ping"))
	assert.NoError(t, err)

	log.Println(resp)
}
