package core

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPingConnections(t *testing.T) {
	peerPool := NewPeerPool(log.Default(), 0, 0)
	defer peerPool.Close()
	assert.Equal(t, 0, peerPool.NumberOfPeers())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m1 := NewMockSender(ctrl)
	m1.EXPECT().SendIsAlive().Return(errors.New("not alive")).AnyTimes()

	m2 := NewMockSender(ctrl)
	m2.EXPECT().SendIsAlive().Return(nil).AnyTimes()

	senders := []Sender{m1, m1, m2, m2}
	for i, s := range senders {
		p := Peer{
			Sender: s,
			addr:   Addr{IP: "127.0.0.1", Port: "809" + fmt.Sprint(i)},
		}
		peerPool.Add(p)
	}

	assert.Equal(t, len(senders), peerPool.NumberOfPeers())

	peerPool.pingConnections()

	assert.Equal(t, 2, peerPool.NumberOfPeers(), "should be equal to the active peer number")
}
