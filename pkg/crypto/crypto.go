package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
)

var pubCurve = elliptic.P256()

func init() {
	// https://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface
	gob.Register(pubCurve)
}

const HashLen = 32

type HashValue = [HashLen]byte

type SecretKey struct {
	sk *ecdsa.PrivateKey
}

type PublicKey struct {
	PubKey ecdsa.PublicKey
}

func (p PublicKey) IsValid() bool {
	return p.PubKey.X != nil &&
		p.PubKey.Y != nil &&
		p.PubKey.Curve != nil &&
		p.PubKey.IsOnCurve(p.PubKey.X, p.PubKey.Y)
}

func Verify(pk PublicKey, signedMsg []byte, r, s *big.Int) bool {
	return ecdsa.Verify(&pk.PubKey, signedMsg, r, s)
}

func NewSecretKey() (SecretKey, error) {
	// this generates a public & private key pair
	sk, err := ecdsa.GenerateKey(pubCurve, rand.Reader)
	if err != nil {
		return SecretKey{}, err
	}

	return SecretKey{sk: sk}, nil
}

func (sk SecretKey) PublicKey() PublicKey {
	return PublicKey{sk.sk.PublicKey}
}

func (sk SecretKey) Sign(message []byte) (*big.Int, *big.Int, error) {
	return ecdsa.Sign(rand.Reader, sk.sk, message)
}

func Hash(message []byte) (HashValue, error) {
	var buf HashValue
	h := sha256.New()
	if _, err := h.Write(message); err != nil {
		return buf, fmt.Errorf("on writing to hash.Hash: %w", err)
	}

	copy(buf[:], h.Sum(nil)[:HashLen])
	return buf, nil
}
