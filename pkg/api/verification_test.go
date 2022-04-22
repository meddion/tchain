package api

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/meddion/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(42)
}

func TestTransactionValidation(t *testing.T) {
	msg := TxData(`message to be hashed`)

	checksum, err := crypto.Hash256(msg[:])
	assert.NoError(t, err, "on hashing a message")

	signer, err := crypto.NewSigner()
	assert.NoError(t, err, "on creating a signer")
	sig, err := signer.Sign(checksum[:])
	assert.NoError(t, err, "on signing a message")

	var invalidSum crypto.HashValue
	copy(invalidSum[5:], checksum[:len(checksum)-5])

	sig2, err := signer.Sign(invalidSum[:])
	assert.NoError(t, err, "on signing a message")

	testTable := []struct {
		tx     Transaction
		target error
	}{
		{Transaction{}, ErrEmptyTxData},
		{Transaction{Data: TxData{}, Hash: checksum}, ErrEmptyTxData},
		{Transaction{Sig: sig, Hash: checksum}, ErrEmptyTxData},
		{Transaction{Data: msg, Hash: checksum}, ErrInvalidSignature},
		{Transaction{Sig: sig2, Data: msg, Hash: checksum}, ErrInvalidSignature},
		{Transaction{Sig: sig, Data: msg}, ErrInvalidChecksum},
		{Transaction{Sig: sig2, Data: msg, Hash: invalidSum}, ErrInvalidChecksum},
		// Success
		{Transaction{Sig: sig, Data: msg, Hash: checksum}, nil},
	}

	for i, testCase := range testTable {
		assert.Equal(t, testCase.target, VerifyTransaction(testCase.tx), "table entry #%d", i)
	}
}

func TestBlockchainValidation(t *testing.T) {
	blocks, err := genRandBlockchain(10)
	assert.NoError(t, err, "on generating blockchain")
	_, genesisBlock := getGenesisPair()

	assert.Equal(t, genesisBlock, blocks[0], "on checking genesis block")

	for i := 1; i < len(blocks); i++ {
		assert.NoError(t, VerifyBlock(blocks[i], blocks[i-1]), "on verifying block")
	}
}

func genRandBlockchain(num int) ([]Block, error) {
	blocks := make([]Block, num)
	_, genesisBlock := getGenesisPair()
	blocks[0] = genesisBlock

	for i := 1; i < len(blocks); i++ {
		hb, err := blocks[i-1].Header.Bytes()
		if err != nil {
			return nil, err
		}

		prevBlockHash, err := crypto.Hash256(hb)
		if err != nil {
			return nil, err
		}

		txs, err := genRandTransactions(25)
		if err != nil {
			return nil, err
		}

		mroot, err := crypto.GenMerkleRoot(txs)
		if err != nil {
			return nil, err
		}

		blocks[i] = Block{
			Header: Header{
				Version:       1,
				Timestamp:     time.Now().Add(time.Second).Unix(),
				PrevBlockHash: prevBlockHash,
				MerkleRoot:    mroot,
				Nonce:         Nonce(i),
			},
			Body: txs,
		}
	}

	return blocks, nil
}

func genRandTransactions(num int) ([]Transaction, error) {
	signer, err := crypto.NewSigner()
	if err != nil {
		return nil, fmt.Errorf("on creating a signer: %w", err)
	}

	txs := make([]Transaction, num)
	for i := 0; i < len(txs); i++ {
		msg := make(TxData, TxBodySizeLimit)
		if _, err := rand.Read(msg); err != nil {
			return nil, fmt.Errorf("on writing a random byte sequence: %w", err)
		}

		hashed, err := crypto.Hash256(msg[:])
		if err != nil {
			return nil, fmt.Errorf("on hashing a message: %w", err)
		}

		sig, err := signer.Sign(hashed[:])
		if err != nil {
			return nil, fmt.Errorf("on signing a message: %w", err)
		}

		txs[i].Data = msg
		txs[i].Hash = hashed
		txs[i].Sig = sig
	}

	return txs, nil
}
