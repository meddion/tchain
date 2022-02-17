package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerkleRoot(t *testing.T) {
	testTable := []struct {
		input    [][]byte
		expected HashValue
	}{
		{
			[][]byte{
				[]byte("0xdead00000"),
				[]byte("0x32fgd2300"),
				[]byte("0x444234"),
			},
			[32]byte{0, 0, 0},
		},
	}

	for _, testCase := range testTable {
		out, err := MerkleRoot(Hash, testCase.input)
		assert.NoError(t, err, "on generating a merkle root hash")
		assert.Equal(t, testCase.expected, out, "on comparing merkle root hashes")
	}
}
