package app

import (
	"github.com/boltdb/bolt"
)

// Db is a single instance of dbHelper
var Db *dbHelper

type dbHelper struct {
	dbPath string
}

func initDb(dbPath string) error {
	// check can open?
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	Db = &dbHelper{dbPath: dbPath}
	return nil
}

func (this *dbHelper) Get(bucket string, key string) ([]byte, error) {
	db, err := bolt.Open(this.dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var value []byte
	err = db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucket))
		if bk == nil {
			return ErrorBucketNotFound
		}
		buf := bk.Get([]byte(key))
		if buf == nil {
			return nil
		}
		value = make([]byte, len(buf))
		copy(value, buf)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (this *dbHelper) Put(bucket string, key string, value []byte) error {
	db, err := bolt.Open(this.dbPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		err = bk.Put([]byte(key), value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (this *dbHelper) Delete(bucket string, key string) error {
	db, err := bolt.Open(this.dbPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucket))
		if bk == nil {
			return ErrorBucketNotFound
		}
		err = bk.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}