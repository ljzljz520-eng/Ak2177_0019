package store

import (
	"sort"

	"inventorychain/internal/model"
)

var LedgerBucket = []byte("ledger")

func (s *Store) SaveLedgerEntry(entry model.LedgerEntry) error {
	return s.put(LedgerBucket, entry.ID, entry)
}

func (s *Store) GetLedgerEntry(id string) (model.LedgerEntry, error) {
	var entry model.LedgerEntry
	err := s.get(LedgerBucket, id, &entry)
	return entry, err
}

func (s *Store) ListLedgerEntries(recordID string) ([]model.LedgerEntry, error) {
	entries := make([]model.LedgerEntry, 0)
	err := s.list(LedgerBucket, func(data []byte) error {
		var entry model.LedgerEntry
		if err := unmarshal(data, &entry); err != nil {
			return err
		}
		if recordID == "" || entry.RecordID == recordID {
			entries = append(entries, entry)
		}
		return nil
	})
	sort.Slice(entries, func(i, j int) bool { return entries[i].Sequence < entries[j].Sequence })
	return entries, err
}

func (s *Store) DeleteLedgerEntry(id string) error {
	return s.delete(LedgerBucket, id)
}
