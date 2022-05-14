package core

import (
	"fmt"
	"time"

	"github.com/meddion/pkg/crypto"
)

var (
	_genesisBlock     Block
	_genesisBlockHash crypto.HashValue
)

func getGenesisPair() (crypto.HashValue, Block) {
	if _genesisBlock.Nonce != 0 {
		return _genesisBlockHash, _genesisBlock
	}

	var err error
	if _genesisBlock, err = genGenesisBlock(); err != nil {
		panic(fmt.Sprintf("on creating a genesis block: %v", err))
	}

	if _genesisBlockHash, err = _genesisBlock.Header.Checksum(); err != nil {
		panic(fmt.Sprintf("on hashing a genesis block: %v", err))
	}

	return _genesisBlockHash, _genesisBlock
}

func genGenesisBlock() (Block, error) {
	ghash, err := crypto.Hash256([]byte("genesis"))
	if err != nil {
		return Block{}, err
	}

	mroot, err := crypto.GenMerkleRoot([]Transaction{})
	if err != nil {
		return Block{}, err
	}

	header := Header{
		Version:       1,
		Timestamp:     time.Date(2021, time.February, 24, 0, 0, 0, 0, time.UTC).Unix(),
		PrevBlockHash: ghash,
		MerkleRoot:    mroot,
		Nonce:         2068160, // difficalty = 21
	}
	// nonce, err := genPowNonce(header, getPowTarget())
	// if err != nil {
	// 	return Block{}, err

	// }
	// header.Nonce = nonce

	return Block{
		Header: header,
		Body:   Body{},
	}, nil
}
