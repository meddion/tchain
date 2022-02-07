package tchain

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const path = "./testdata/db"

func TestDB(t *testing.T) {
	assert.NoError(t, os.RemoveAll(path))

	db, err := NewDB(path)
	assert.NoError(t, err, "on creating DB instance")
	//
	// batch := new(leveldb.Batch)
	// batch.Put([]byte("foo"), []byte("foo"))
	// batch.Put([]byte("foo-baz"), []byte("bar"))
	// batch.Put([]byte("foo-bar"), []byte("foo-bar"))
	// batch.Put([]byte("foo-baz"), []byte("zzz"))
	// // batch.Delete([]byte("foo-baz"))
	// assert.NoError(t, db.Write(batch, nil))
	//
	// iter := db.NewIterator(util.BytesPrefix([]byte("foo-")), nil)
	// defer iter.Release()
	//
	// for iter.Next() {
	//     log.Printf("%s\n", iter.Value())
	// }
	// assert.NoError(t, iter.Error())

	defer func() {
		assert.NoError(t, db.Close())
	}()
}
