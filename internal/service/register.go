package service

import (
	"fmt"

	"inventorychain/internal/flow015"
	"inventorychain/internal/model"
	"inventorychain/internal/review"
)

func (s *Service) Register(record model.Record, actor string) (model.Record, error) {
	if err := s.requireStore(); err != nil {
		return model.Record{}, err
	}
	if record.ID == "" {
		record.ID = model.RecordID(record.Warehouse, record.Cycle)
	}
	if record.Status == "" {
		record.Status = model.StatusDraft
	}
	record.CreatedBy = actor
	record.UpdatedBy = actor
	record.Version = 1
	record.CreatedAt = s.clock.Now()
	record.UpdatedAt = record.CreatedAt
	review.NormalizeLines(&record)
	if err := record.Validate(); err != nil {
		return model.Record{}, err
	}
	assignment, err := flow015.BuildAssignment(record)
	if err != nil {
		return model.Record{}, err
	}
	if !assignment.OrderOK {
		return model.Record{}, fmt.Errorf("slot order is invalid")
	}
	if _, err := s.store.GetRecord(record.ID); err == nil {
		return model.Record{}, fmt.Errorf("record %s already exists", record.ID)
	}
	if err := s.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, s.appendEvent(record.ID, "register", actor, "record registered")
}

func (s *Service) Submit(recordID, actor string) (model.Record, error) {
	return s.transition(recordID, "submit", actor, "submitted for review")
}

func (s *Service) Approve(recordID, actor string) (model.Record, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if err := review.RequireApproval(record); err != nil {
		return model.Record{}, err
	}
	return s.transitionWithRecord(record, "approve", actor, "approved after review")
}

func (s *Service) Reject(recordID, actor, reason string) (model.Record, error) {
	return s.transition(recordID, "reject", actor, model.CleanReason(reason))
}

func (s *Service) Archive(recordID, actor string) (model.Record, error) {
	return s.transition(recordID, "archive", actor, "archived")
}

func (s *Service) transition(recordID, action, actor, detail string) (model.Record, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return model.Record{}, err
	}
	return s.transitionWithRecord(record, action, actor, detail)
}

func (s *Service) transitionWithRecord(record model.Record, action, actor, detail string) (model.Record, error) {
	next, ok := model.NextStatus(record.Status, action)
	if !ok {
		return model.Record{}, fmt.Errorf("cannot %s record in status %s", action, record.Status)
	}
	record.Status = next
	record.UpdatedBy = actor
	record.UpdatedAt = s.clock.Now()
	record.Version++
	if err := s.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, s.appendEvent(record.ID, action, actor, detail)
}

func (s *Service) appendEvent(recordID, action, actor, detail string) error {
	events, err := s.store.ListEvents(recordID)
	if err != nil {
		return err
	}
	event := model.AuditEvent{RecordID: recordID, Action: action, Actor: actor, Detail: detail, Sequence: len(events) + 1, Timestamp: s.clock.Now()}
	event.ID = model.EventID(recordID, event.Sequence)
	return s.store.SaveEvent(event)
}
