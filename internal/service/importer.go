package service

import (
	"fmt"

	"inventorychain/internal/model"
	"inventorychain/internal/report"
)

type ImportRow struct {
	Warehouse string
	Cycle     string
	SKU       string
	Expected  int
	Observed  int
	Reason    string
	Slot      string
}

type ImportReport struct {
	Accepted []model.Record
	Rejected []string
	Summary  report.Summary
}

func (s *Service) Import(rows []ImportRow, actor string) (ImportReport, error) {
	if err := s.requireStore(); err != nil {
		return ImportReport{}, err
	}
	grouped := make(map[string]model.Record)
	order := make([]string, 0)
	result := ImportReport{Accepted: []model.Record{}, Rejected: []string{}}
	for _, row := range rows {
		if row.Warehouse == "" || row.Cycle == "" || row.SKU == "" {
			result.Rejected = append(result.Rejected, "row has missing identity")
			continue
		}
		id := model.RecordID(row.Warehouse, row.Cycle)
		record, exists := grouped[id]
		if !exists {
			record = model.Record{ID: id, Warehouse: row.Warehouse, Cycle: row.Cycle, Status: model.StatusDraft, Slots: []model.TimeSlot{{Name: row.Slot, Sequence: 1, Owner: actor}}, Lines: []model.DiscrepancyLine{}}
			grouped[id] = record
			order = append(order, id)
		}
		record.Lines = append(record.Lines, model.DiscrepancyLine{SKU: row.SKU, Expected: row.Expected, Observed: row.Observed, Reason: row.Reason, Slot: row.Slot})
		grouped[id] = record
	}
	for _, id := range order {
		record, err := s.Register(grouped[id], actor)
		if err != nil {
			result.Rejected = append(result.Rejected, fmt.Sprintf("%s: %v", id, err))
			continue
		}
		result.Accepted = append(result.Accepted, record)
	}
	result.Summary = report.BuildSummary(result.Accepted)
	return result, nil
}

func (s *Service) Attach(recordID, actor, name, kind, content string) (model.Attachment, error) {
	if err := s.requireStore(); err != nil {
		return model.Attachment{}, err
	}
	if _, err := s.Get(recordID); err != nil {
		return model.Attachment{}, err
	}
	attachment := model.Attachment{ID: model.AttachmentID(recordID, name), RecordID: recordID, Name: name, Kind: kind, Content: content, Checksum: fmt.Sprintf("%x", len(content))}
	if err := s.store.SaveAttachment(attachment); err != nil {
		return model.Attachment{}, err
	}
	if err := s.appendEvent(recordID, "attach", actor, "attachment added"); err != nil {
		return model.Attachment{}, err
	}
	return attachment, nil
}
