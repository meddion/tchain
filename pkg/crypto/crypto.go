package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"math/big"
)

var _pubCurve = elliptic.P256()

func init() {
	// https://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface
	gob.Register(_pubCurve)
	gob.Register(sigECDSA{})
}

type SignerECDSA struct {
	sk *ecdsa.PrivateKey
}

func NewSignerECDSA() (SignerECDSA, error) {
	// this generates a public & private key pair
	sk, err := ecdsa.GenerateKey(_pubCurve, rand.Reader)
	if err != nil {
		return SignerECDSA{}, err
	}

	return SignerECDSA{sk: sk}, nil
}

func (sk SignerECDSA) Sign(message []byte) (sigECDSA, error) {
	r, s, err := ecdsa.Sign(rand.Reader, sk.sk, message)
	return sigECDSA{PK: sk.sk.PublicKey, R: r, S: s}, err
}

type sigECDSA struct {
	PK   ecdsa.PublicKey
	R, S *big.Int
}

func (sig sigECDSA) Verify(signedMsg []byte) bool {
	return sig.isValidPubKey() && ecdsa.Verify(&sig.PK, signedMsg, sig.R, sig.S)
}

func (sig sigECDSA) isValidPubKey() bool {
	return sig.PK.X != nil &&
		sig.PK.Y != nil &&
		sig.PK.Curve != nil &&
		sig.PK.IsOnCurve(sig.PK.X, sig.PK.Y)
}
