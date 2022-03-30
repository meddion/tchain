package db

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipLists_Search(t *testing.T) {
	table := []struct{ k, v []byte }{
		{[]byte(`key1`), []byte(`val1`)},
		{[]byte(`key2`), []byte(`val2`)},
		{[]byte(`key3`), []byte(`val3`)},
		{[]byte(`key4`), []byte(`val4`)},
	}

	skipLists := NewSkipLists()
	for _, e := range table {
		n := skipLists.Insert(e.k, e.v)
		assert.True(t, bytes.Equal(n.value, e.v))
	}
	// assert.Equal(t, skipLists.height, len(table))
}
