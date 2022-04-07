package api

import (
	"context"
	"log"
	"math/rand"
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
	s.NoError(s.db.Store(GenesisBlockHash, GenesisBlock), "on storing a genesis block")

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

	var msg TxData
	copy(msg[:], []byte(`This is not enough!`))

	hashed, err := crypto.Hash256(msg[:])
	s.NoError(err, "on hashing a message")

	sig, err := s.signer.Sign(hashed[:])
	s.NoError(err, "on signing a message")

	testTable := []struct {
		tx  Transaction
		err error
	}{
		// error
		{Transaction{Data: msg, Hash: hashed}, ErrInvalidTransaction},
		{Transaction{Data: TxData{}, Hash: hashed}, ErrInvalidTransaction},
		{Transaction{Sig: sig, Data: msg}, ErrInvalidTransaction},
		// success
		{Transaction{Sig: sig, Data: msg, Hash: hashed}, nil},
	}

	for _, testCase := range testTable {
		_, err = c.SendTransaction(TransactionReq{Transaction: testCase.tx})
		if testCase.err == nil {
			s.NoError(err)
		} else {
			s.ErrorAs(err, &testCase.err, "on executing SendTransaction()")
		}
	}
}

func (s *senderReceiverSuite) TestBlocks() {
	c := s.peerPool.Peers()[0]

	txs := make([]Transaction, 25)
	for i := 0; i < len(txs); i++ {
		var msg TxData
		_, err := rand.Read(msg[:])
		s.NoError(err, "on writing a random byte sequence")

		hashed, err := crypto.Hash256(msg[:])
		s.NoError(err, "on hashing a message")

		sig, err := s.signer.Sign(hashed[:])
		s.NoError(err, "on signing a message")

		txs[i].Data = msg
		txs[i].Hash = hashed
		txs[i].Sig = sig
	}

	txBytes := make([][]byte, 0, len(txs))
	for _, tx := range txs {
		btx, err := tx.Bytes()
		s.NoError(err)
		txBytes = append(txBytes, btx)
	}

	merkleRoot, err := crypto.GenMerkleRoot(txBytes)
	s.NoError(err, "on generating merkle root")

	prevBlockPair := s.db.LastCommited()

	var body Body
	copy(body[:], txs)

	newBlockReq := BlockReq{
		Block{
			Header: Header{
				version:       1,
				Timestamp:     time.Now().Unix(),
				PrevBlockHash: prevBlockPair.hash,
				MerkleRoot:    merkleRoot,
				Nonce:         2,
			},
			Body: body,
		},
	}

	s.NoError(c.SendBlock(newBlockReq), "on commiting a new block")
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(senderReceiverSuite))
}
