package main

import (
	bolt "go.etcd.io/bbolt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
)

type Record struct {
	Key  string
	Type string
	URL  string
}

// generateKey generates a key of specified length from a predefined list of characters
func generateKey(size int) string {
	chars := "abcdefghijkmnopqrstuvwxyz23456789ABCDEFGHJKMNPQRSTUVWXYZ"

	rand.Seed(time.Now().UTC().UnixNano())

	out := ""
	for i := 0; len(out) < size; i++ {
		out = out + string(chars[rand.Intn(len(chars))])
	}

	return out
}

// newKey writes a new key to the database
func saveRecord(rec Record) error {
	data, err := bson.Marshal(rec)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Records"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(rec.Key), []byte(data))
		return err
	})

	return err
}

// getKey gets the record from the database identified by the given key
func getKey(key string) (Record, error) {

	var data []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Records"))
		data = b.Get([]byte(key))
		return nil
	})

	rec := Record{}
	err = bson.Unmarshal(data, &rec)
	return rec, err
}
