package core

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

type mockPeerPool struct {
	PeerPool
}

func (m mockPeerPool) Peers() []Peer {
	addr := Addr{_testAddr, _testPort}
	s, err := NewSender(addr)
	if err != nil {
		log.Fatalf("on creating a test sender: %v", err)
	}

	return []Peer{{Sender: s, addr: addr}}
}

func (m mockPeerPool) NumberOfPeers() int {
	return 1
}

func (m mockPeerPool) SendToPeers(_ func(Peer) error) <-chan error {
	c := make(chan error)
	close(c)

	return c
}

func (m mockPeerPool) Close() error {
	return nil
}

type senderReceiverSuite struct {
	suite.Suite

	signer   crypto.SignerECDSA
	peerPool PeerPool
	blkchain *Blockchain
	serv     *Server
}

func (s *senderReceiverSuite) SetupSuite() {
	b, err := bolt.Open(_dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		s.NoError(err, "on opening a bolt conn")
	}

	db, err := NewBlockRepo(b)
	s.NoError(err, "on creating a block repo")

	logger := log.Default()
	s.blkchain, err = NewBlockchain(db, logger)

	s.NoError(err, "on creating the Blockchain instance")

	s.peerPool = mockPeerPool{}

	rcv := NewReceiverRPC(s.blkchain, s.peerPool, logger)

	s.signer, err = crypto.NewSignerECDSA()
	s.NoError(err, "on creating a signer")

	s.serv, err = NewServer(rcv)
	s.NoError(err, "on creating a server")

	go func() {
		s.Equal(http.ErrServerClosed, s.serv.Start(_testAddr, _testPort), "on closing a server")
	}()
	<-time.After(time.Millisecond * 500)
}

func (s *senderReceiverSuite) TearDownSuite() {
	s.NoError(s.serv.Close(context.Background()), "on closing database")
	s.NoError(s.peerPool.Close(), "on closing the peer pool")
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

	h := Header{
		Version:       1,
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: _genesisBlockHash,
		MerkleRoot:    merkleRoot,
		Difficulty:    Difficulty(15),
	}

	nonce, err := h.Difficulty.GenNonce(h)
	s.NoError(err, "on generating nonce for header")
	h.Nonce = nonce

	newBlockReq := BlockReq{
		Block{
			Header: h,
			Body:   txs,
		},
	}

	s.NoError(c.SendBlock(newBlockReq), "on commiting a new block")
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(senderReceiverSuite))
}
