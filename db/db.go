package db

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/pindamonhangaba/urlshorts/service"
	bolt "go.etcd.io/bbolt"
)

const (
	// BucketName is the name of the bucket where URL data is stored
	BucketName = "urls"
)

// DB represents a database instance
type DB struct {
	db *bolt.DB
}

// NewDB creates a new database instance
func NewDB(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Create the bucket if it doesn't exist
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		return err
	}); err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// Close closes the database
func (d *DB) Close() error {
	return d.db.Close()
}

// SaveURL saves a URL to the database
func (d *DB) SaveURL(url *service.URL) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}

		// Marshal the URL struct to JSON
		urlJSON, err := json.Marshal(url)
		if err != nil {
			return err
		}

		// Save the URL to the database
		return bucket.Put([]byte(url.Code), urlJSON)
	})
}

// GetURL retrieves a URL from the database by its code
func (d *DB) GetURL(code string) (*service.URL, error) {
	var url service.URL

	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}

		data := bucket.Get([]byte(code))
		if data == nil {
			return errors.New("URL not found")
		}

		return json.Unmarshal(data, &url)
	})

	if err != nil {
		return nil, err
	}

	return &url, nil
}

// ListURLs returns all URLs in the database
func (d *DB) ListURLs() ([]*service.URL, error) {
	var urls []*service.URL

	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var url service.URL
			if err := json.Unmarshal(v, &url); err != nil {
				return err
			}
			urls = append(urls, &url)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return urls, nil
}

// DeleteURL deletes a URL from the database by its code
func (d *DB) DeleteURL(code string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}

		return bucket.Delete([]byte(code))
	})
}
