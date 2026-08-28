package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"go.etcd.io/bbolt"
	"inventorychain/internal/model"
)

var (
	RecordsBucket     = []byte("records")
	EventsBucket      = []byte("events")
	WorkflowsBucket   = []byte("workflows")
	AttachmentsBucket = []byte("attachments")
)

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db}
	if err := s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{RecordsBucket, EventsBucket, WorkflowsBucket, AttachmentsBucket, LedgerBucket} {
			if _, createErr := tx.CreateBucketIfNotExists(bucket); createErr != nil {
				return createErr
			}
		}
		return nil
	}); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize buckets: %w", err)
	}
	return s, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func marshal(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode value: %w", err)
	}
	return data, nil
}

func unmarshal(data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode value: %w", err)
	}
	return nil
}

func (s *Store) put(bucket []byte, key string, value any) error {
	data, err := marshal(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Put([]byte(key), data)
	})
}

func (s *Store) get(bucket []byte, key string, target any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucket).Get([]byte(key))
		if value == nil {
			return fmt.Errorf("%s %q not found", string(bucket), key)
		}
		return unmarshal(value, target)
	})
}

func (s *Store) delete(bucket []byte, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Delete([]byte(key))
	})
}

func (s *Store) list(bucket []byte, decode func([]byte) error) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			return decode(value)
		})
	})
}

func (s *Store) SaveRecord(record model.Record) error {
	return s.put(RecordsBucket, record.ID, record)
}

func (s *Store) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := s.get(RecordsBucket, id, &record)
	return record, err
}

func (s *Store) DeleteRecord(id string) error {
	return s.delete(RecordsBucket, id)
}
