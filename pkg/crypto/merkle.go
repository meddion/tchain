package crypto

import (
	"crypto/sha256"
	"fmt"
)

const HashLen uint8 = 32

type (
	HashValue = [HashLen]byte
	HashFunc  func([]byte) (HashValue, error)
)

var DefaultHashFunc HashFunc = Hash256

func Hash256(message []byte) (HashValue, error) {
	var buf HashValue
	h := sha256.New()
	if _, err := h.Write(message); err != nil {
		return buf, fmt.Errorf("on writing to hash.Hash: %w", err)
	}

	copy(buf[:], h.Sum(nil)[:HashLen])
	return buf, nil
}

func GenMerkleRoot(values [][]byte) (HashValue, error) {
	switch len(values) {
	case 0:
		return DefaultHashFunc([]byte{})
	case 1:
		return DefaultHashFunc(values[0])
	}

	N := len(values)
	if len(values)%2 != 0 {
		values = append(values, values[len(values)-1])
		N++
	}
	hashes := make([]HashValue, N)

	for i, v := range values {
		hash, err := DefaultHashFunc(v)
		if err != nil {
			return HashValue{}, nil
		}
		hashes[i] = hash
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes[len(hashes)] = hashes[len(hashes)-1]
			hashes = hashes[:len(hashes)+1]
		}

		for i, j := 0, 0; i < len(hashes); i, j = i+2, j+1 {
			buf := make([]byte, HashLen*2)
			copy(buf[:HashLen], hashes[i][:])
			copy(buf[HashLen:], hashes[i+1][:])

			hash, err := DefaultHashFunc(buf)
			if err != nil {
				return HashValue{}, nil
			}
			hashes[j] = hash
		}

		hashes = hashes[:len(hashes)/2]
	}

	return hashes[0], nil
}
