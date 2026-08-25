package store

import (
	"fmt"
	"sort"

	"go.etcd.io/bbolt"
	"inventorychain/internal/model"
)

func (s *Store) ListRecords() ([]model.Record, error) {
	records := make([]model.Record, 0)
	err := s.list(RecordsBucket, func(data []byte) error {
		var record model.Record
		if err := unmarshal(data, &record); err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	sort.Slice(records, func(i, j int) bool {
		if records[i].Warehouse == records[j].Warehouse {
			return records[i].Cycle < records[j].Cycle
		}
		return records[i].Warehouse < records[j].Warehouse
	})
	return records, err
}

func (s *Store) UpdateRecord(id string, update func(*model.Record) error) (model.Record, error) {
	var result model.Record
	err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(RecordsBucket)
		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("records %q not found", id)
		}
		if err := unmarshal(data, &result); err != nil {
			return err
		}
		if err := update(&result); err != nil {
			return err
		}
		encoded, err := marshal(result)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(id), encoded)
	})
	return result, err
}
