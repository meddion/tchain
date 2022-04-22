package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/meddion/pkg/crypto"
	"github.com/stretchr/testify/suite"
)

const (
	_dbFile   = "_test_db_file_"
	_testAddr = ""
	_testPort = "2022"
)

type mockPeerPool struct{}

func (m mockPeerPool) Peers() []Sender {
	s, err := NewSender(_testAddr, _testPort)
	if err != nil {
		log.Fatalf("on creating a test sender: %v", err)
	}

	return []Sender{s}
}

type senderReceiverSuite struct {
	suite.Suite

	signer   crypto.Signer
	peerPool PeerPool
	db       *BlockRepo
	serv     *Server
}

func (s *senderReceiverSuite) SetupSuite() {
	b, err := bolt.Open(_dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		s.NoError(err, "on opening a bolt conn")
	}

	s.db, err = NewBlockRepo(b)
	s.NoError(err, "on creating a block repo")

	// Setup
	s.NoError(s.db.Store(getGenesisPair()), "on storing a genesis block")

	s.peerPool = mockPeerPool{}
	rcv := NewReceiverRPC(s.peerPool, s.db, log.Default())

	s.signer, err = crypto.NewSigner()
	s.NoError(err, "on creating a signer")

	s.serv, err = NewServer(rcv)
	s.NoError(err, "on creating a server")

	go func() { s.Equal(http.ErrServerClosed, s.serv.Start(_testAddr, _testPort), "on closing a server") }()
	<-time.After(time.Millisecond * 500)
}

func (s *senderReceiverSuite) TearDownSuite() {
	s.NoError(s.serv.Close(context.Background()), "on closing database")
	s.NoError(os.Remove(_dbFile), "on removing a file")
}

func (s *senderReceiverSuite) TestTransactions() {
	c := s.peerPool.Peers()[0]

	msg := TxData(`This is not enough!`)

	hashed, err := crypto.Hash256(msg[:])
	s.NoError(err, "on hashing a message")

	sig, err := s.signer.Sign(hashed[:])
	s.NoError(err, "on signing a message")

	testTable := []struct {
		tx  Transaction
		err error
	}{
		// error
		{Transaction{Data: msg, Hash: hashed}, ErrInvalidSignature},
		{Transaction{Data: TxData{}, Hash: hashed}, ErrEmptyTxData},
		{Transaction{Sig: sig, Data: msg}, ErrInvalidChecksum},
		// success
		{Transaction{Sig: sig, Data: msg, Hash: hashed}, nil},
	}

	for _, testCase := range testTable {
		_, err = c.SendTransaction(TransactionReq{Transaction: testCase.tx})
		if testCase.err == nil {
			s.NoError(err)
		} else {
			s.ErrorAs(err, &testCase.err)
		}
	}
}
func (s *senderReceiverSuite) TestBlocks() {
	c := s.peerPool.Peers()[0]

	txs, err := genRandTransactions(25)
	s.NoError(err, "generating random txs")

	merkleRoot, err := crypto.GenMerkleRoot(txs)
	s.NoError(err, "on generating merkle root")

	prevBlockPair := s.db.LastCommited()

	newBlockReq := BlockReq{
		Block{
			Header: Header{
				Version:       1,
				Timestamp:     time.Now().Unix(),
				PrevBlockHash: prevBlockPair.hash,
				MerkleRoot:    merkleRoot,
				Nonce:         2,
			},
			Body: txs,
		},
	}

	s.NoError(c.SendBlock(newBlockReq), "on commiting a new block")
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(senderReceiverSuite))
}
