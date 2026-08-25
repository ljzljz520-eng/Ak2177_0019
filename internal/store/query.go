package store

import (
	"strings"

	"inventorychain/internal/model"
)

func (s *Store) FindRecordsByStatus(status model.RecordStatus) ([]model.Record, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]model.Record, 0)
	for _, record := range records {
		if record.Status == status {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (s *Store) CountRecordsByWarehouse() (map[string]int, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, record := range records {
		key := strings.TrimSpace(record.Warehouse)
		counts[key]++
	}
	return counts, nil
}

func (s *Store) RemoveArchived(ids []string) (int, error) {
	removed := 0
	for _, id := range ids {
		record, err := s.GetRecord(id)
		if err != nil {
			continue
		}
		if record.Status == model.StatusArchived {
			if err := s.DeleteRecord(id); err != nil {
				return removed, err
			}
			removed++
		}
	}
	return removed, nil
}

func (s *Store) ReplaceRecord(record model.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.SaveRecord(record)
}
