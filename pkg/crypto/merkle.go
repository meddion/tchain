package crypto

type HashFunc func([]byte) HashValue

func MerkleRoot(hashFunc HashFunc, values [][]byte) HashValue {
}
