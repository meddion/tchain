package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: add more inteligent tests
type bytes []byte

func (b bytes) Bytes() ([]byte, error) {
	return b, nil
}

func TestMerkleRoot(t *testing.T) {

	testTable := []struct {
		input    []bytes
		expected HashValue
	}{
		{
			[]bytes{
				bytes("0xdead00000"),
				bytes("0x32fgd2300"),
				bytes("0x444234"),
			},
			[32]byte{0x1d, 0x1f, 0xc3, 0x5b, 0xf2, 0x5, 0x9e, 0xb5, 0x9d, 0x53,
				0xec, 0xb5, 0xa6, 0x63, 0x70, 0x7c, 0x45, 0xaf, 0x41, 0x3b, 0x7,
				0x86, 0x12, 0x6f, 0x5c, 0x6a, 0xe2, 0xbf, 0xda, 0x22, 0xa9, 0xe2,
			},
		},
	}

	for _, testCase := range testTable {
		out, err := GenMerkleRoot(testCase.input)
		assert.NoError(t, err, "on generating a merkle root hash")
		assert.Equal(t, testCase.expected, out, "on comparing merkle root hashes")
	}
}
