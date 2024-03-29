package crypto

import (
	"crypto/sha256"
	"fmt"
)

const HashLen uint8 = 32

type HashValue = [HashLen]byte

var (
	ZeroHashValue    HashValue
	_defaultHashFunc func([]byte) (HashValue, error) = Hash256
)

func Hash256(message []byte) (HashValue, error) {
	var buf HashValue
	h := sha256.New()
	if _, err := h.Write(message); err != nil {
		return buf, fmt.Errorf("on writing to hash.Hash: %w", err)
	}

	copy(buf[:], h.Sum(nil)[:HashLen])
	return buf, nil
}

type bytesConverter interface {
	Bytes() ([]byte, error)
}

func GenMerkleRoot[T bytesConverter](values []T) (HashValue, error) {
	switch len(values) {
	case 0:
		return _defaultHashFunc([]byte{})
	case 1:
		v, err := values[0].Bytes()
		if err != nil {
			return HashValue{}, err
		}

		return _defaultHashFunc(v)
	}

	hashes := make([]HashValue, len(values))
	for i, v := range values {
		v, err := v.Bytes()
		if err != nil {
			return HashValue{}, err
		}

		hash, err := _defaultHashFunc(v)
		if err != nil {
			return HashValue{}, nil
		}
		hashes[i] = hash
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}

		for i, j := 0, 0; i < len(hashes); i, j = i+2, j+1 {
			buf := make([]byte, HashLen*2)
			copy(buf[:HashLen], hashes[i][:])
			copy(buf[HashLen:], hashes[i+1][:])

			hash, err := _defaultHashFunc(buf)
			if err != nil {
				return HashValue{}, nil
			}
			hashes[j] = hash
		}

		hashes = hashes[:len(hashes)/2]
	}

	return hashes[0], nil
}
