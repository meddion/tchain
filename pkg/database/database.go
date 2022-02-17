package database

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
	*leveldb.DB
}

func NewDB(path string) (DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return DB{}, err
	}
	return DB{
		db,
	}, nil
}

func (db DB) Close() error {
	return db.DB.Close()
}
