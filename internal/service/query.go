package service

import (
	"fmt"

	"inventorychain/internal/model"
	"inventorychain/internal/report"
)

func (s *Service) Get(recordID string) (model.Record, error) {
	if err := s.requireStore(); err != nil {
		return model.Record{}, err
	}
	return s.store.GetRecord(recordID)
}

func (s *Service) Search(filter model.SearchFilter, page, size int) (model.SearchResult, error) {
	if err := s.requireStore(); err != nil {
		return model.SearchResult{}, err
	}
	records, err := s.store.ListRecords()
	if err != nil {
		return model.SearchResult{}, err
	}
	filtered := make([]model.Record, 0, len(records))
	for _, record := range records {
		if filter.Matches(record) {
			filtered = append(filtered, record)
		}
	}
	return model.Paginate(filtered, page, size), nil
}

func (s *Service) UpdateLines(recordID, actor string, lines []model.DiscrepancyLine) (model.Record, error) {
	if err := s.requireStore(); err != nil {
		return model.Record{}, err
	}
	record, err := s.store.UpdateRecord(recordID, func(record *model.Record) error {
		if record.Status == model.StatusArchived {
			return fmt.Errorf("archived record cannot be changed")
		}
		record.Lines = append([]model.DiscrepancyLine(nil), lines...)
		record.UpdatedBy = actor
		record.UpdatedAt = s.clock.Now()
		record.Version++
		return record.Validate()
	})
	if err != nil {
		return model.Record{}, err
	}
	if err := s.appendEvent(recordID, "update", actor, "discrepancy lines updated"); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Publish(recordID, actor string) (model.Record, error) {
	if err := s.requireStore(); err != nil {
		return model.Record{}, err
	}
	record, err := s.store.UpdateRecord(recordID, func(record *model.Record) error {
		if record.Status != model.StatusApproved {
			return fmt.Errorf("only approved records can be published")
		}
		record.Published = true
		record.UpdatedBy = actor
		record.UpdatedAt = s.clock.Now()
		record.Version++
		return nil
	})
	if err != nil {
		return model.Record{}, err
	}
	return record, s.appendEvent(recordID, "publish", actor, "published to warehouse ledger")
}

func (s *Service) ExportCSV(filter model.SearchFilter) (string, error) {
	result, err := s.Search(filter, 1, 10000)
	if err != nil {
		return "", err
	}
	return report.CSV(result.Records), nil
}
