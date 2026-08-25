package service

import (
	"fmt"

	"inventorychain/internal/flow015"
	"inventorychain/internal/model"
	"inventorychain/internal/report"
)

type AdjustmentRequest struct {
	SKU       string
	Quantity  int
	Direction string
	Reason    string
}

func (s *Service) AddAdjustment(recordID, actor string, request AdjustmentRequest) (model.LedgerEntry, error) {
	if err := s.requireStore(); err != nil {
		return model.LedgerEntry{}, err
	}
	if request.SKU == "" || request.Quantity < 0 {
		return model.LedgerEntry{}, fmt.Errorf("adjustment requires sku and non-negative quantity")
	}
	if _, err := s.Get(recordID); err != nil {
		return model.LedgerEntry{}, err
	}
	entries, err := s.store.ListLedgerEntries(recordID)
	if err != nil {
		return model.LedgerEntry{}, err
	}
	entry := model.NewLedgerEntry(recordID, request.SKU, request.Direction, actor, s.clock.Now(), request.Quantity, len(entries)+1, request.Reason)
	if err := s.store.SaveLedgerEntry(entry); err != nil {
		return model.LedgerEntry{}, err
	}
	if err := s.appendEvent(recordID, "adjust", actor, fmt.Sprintf("ledger entry %s", entry.SKU)); err != nil {
		return model.LedgerEntry{}, err
	}
	return entry, nil
}

func (s *Service) Reconcile(recordID string) (model.Reconciliation, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return model.Reconciliation{}, err
	}
	entries, err := s.store.ListLedgerEntries(recordID)
	if err != nil {
		return model.Reconciliation{}, err
	}
	return model.BuildReconciliation(record, entries), nil
}

func (s *Service) SlotReport(recordID string) (flow015.ReconcileResult, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return flow015.ReconcileResult{}, err
	}
	assignment, err := flow015.BuildAssignment(record)
	if err != nil {
		return flow015.ReconcileResult{}, err
	}
	return flow015.ReconcileSlots(record.Slots, assignment.Slots), nil
}

func (s *Service) DetailReport(recordID string) (report.Detail, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return report.Detail{}, err
	}
	events, err := s.store.ListEvents(recordID)
	if err != nil {
		return report.Detail{}, err
	}
	attachments, err := s.store.ListAttachments(recordID)
	if err != nil {
		return report.Detail{}, err
	}
	entries, err := s.store.ListLedgerEntries(recordID)
	if err != nil {
		return report.Detail{}, err
	}
	return report.BuildDetail(record, events, attachments, entries), nil
}

func (s *Service) ReopenWorkflow(recordID, actor string) (model.Workflow, error) {
	if err := s.requireStore(); err != nil {
		return model.Workflow{}, err
	}
	if _, err := s.Get(recordID); err != nil {
		return model.Workflow{}, err
	}
	workflow := model.Workflow{ID: model.WorkflowID(recordID, "inventory-review"), RecordID: recordID, Name: "inventory-review", State: "active", Steps: []string{"register", "review", "confirm", "archive"}, Current: 0, Owner: actor}
	if err := s.store.SaveWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	return workflow, nil
}

func (s *Service) AdvanceWorkflow(recordID, actor string) (model.Workflow, error) {
	workflows, err := s.store.ListWorkflows(recordID)
	if err != nil {
		return model.Workflow{}, err
	}
	if len(workflows) == 0 {
		return model.Workflow{}, fmt.Errorf("workflow for %s not found", recordID)
	}
	workflow := workflows[0]
	if workflow.Current < len(workflow.Steps)-1 {
		workflow.Current++
	}
	if workflow.Current == len(workflow.Steps)-1 {
		workflow.State = "complete"
	}
	workflow.Owner = actor
	if err := s.store.SaveWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	return workflow, nil
}
