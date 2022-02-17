package api

import (
	"log"
	"testing"
	"time"

	"github.com/meddion/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	port = "2022"
	addr = ""
)

type mockSenderPool struct {
	senders func() []Sender
}

func (m mockSenderPool) Senders() []Sender {
	return m.senders()
}

type mockMemPool struct {
	MemPool
}

func (m mockMemPool) Get(_ crypto.HashValue) (Transaction, bool) {
	return Transaction{}, false
}

func TestTransactionRPC(t *testing.T) {
	mockMemPool := mockMemPool{}
	mockSenderPool := mockSenderPool{
		senders: func() []Sender { return []Sender{} },
	}

	rcv := NewReceiverRPC(mockMemPool, mockSenderPool, log.Default())
	s, err := NewServer(rcv)
	assert.NoError(t, err)

	sk, err := crypto.NewSecretKey()
	assert.NoError(t, err)

	go func() { assert.NoError(t, s.Start(addr, port)) }()
	<-time.After(time.Second)

	t.Run("error_transaction", func(t *testing.T) {
		c, err := NewSender(mockSenderPool, addr, port)
		assert.NoError(t, err, "on creating a Sender")

		msg := []byte(`This is not enough!`)
		hashed, err := crypto.Hash(msg)
		assert.NoError(t, err, "on hashing a message")

		r, s, err := sk.Sign(hashed[:])
		assert.NoError(t, err, "on signing a message")

		testTable := []struct {
			tx  Transaction
			err error
		}{
			// error
			{Transaction{Data: msg, Hash: hashed}, InvalidTransactionError},
			{Transaction{Data: make([]byte, 0), Hash: hashed}, InvalidTransactionError},
			{Transaction{PublicKey: sk.PublicKey(), Data: msg, Hash: hashed, R: s, S: r}, InvalidTransactionError},
			// success
			{Transaction{PublicKey: sk.PublicKey(), Data: msg, Hash: hashed, R: r, S: s}, nil},
		}

		for _, testCase := range testTable {
			_, err = c.SendTransaction(TransactionReq{Transaction: testCase.tx})
			if testCase.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &testCase.err, "on executing SendTransaction()")
			}
		}
	})
}
