package store

import (
	"fmt"

	"inventorychain/internal/model"
)

type IntegrityReport struct {
	Records       int
	Events        int
	Workflows     int
	Attachments   int
	LedgerEntries int
	Issues        []string
}

func (s *Store) CheckIntegrity() (IntegrityReport, error) {
	report := IntegrityReport{Issues: []string{}}
	records, err := s.ListRecords()
	if err != nil {
		return report, err
	}
	report.Records = len(records)
	for _, record := range records {
		if err := record.Validate(); err != nil {
			report.Issues = append(report.Issues, fmt.Sprintf("record %s: %v", record.ID, err))
		}
		events, eventErr := s.ListEvents(record.ID)
		if eventErr != nil {
			return report, eventErr
		}
		report.Events += len(events)
		for _, event := range events {
			if event.RecordID != record.ID {
				report.Issues = append(report.Issues, fmt.Sprintf("event %s points to %s", event.ID, event.RecordID))
			}
		}
		attachments, attachmentErr := s.ListAttachments(record.ID)
		if attachmentErr != nil {
			return report, attachmentErr
		}
		report.Attachments += len(attachments)
		entries, ledgerErr := s.ListLedgerEntries(record.ID)
		if ledgerErr != nil {
			return report, ledgerErr
		}
		report.LedgerEntries += len(entries)
		for _, entry := range entries {
			if !entry.Valid() {
				report.Issues = append(report.Issues, fmt.Sprintf("ledger entry %s is invalid", entry.ID))
			}
		}
	}
	workflows, err := s.ListWorkflows("")
	if err != nil {
		return report, err
	}
	report.Workflows = len(workflows)
	for _, workflow := range workflows {
		if workflow.RecordID == "" || len(workflow.Steps) == 0 {
			report.Issues = append(report.Issues, fmt.Sprintf("workflow %s is incomplete", workflow.ID))
		}
	}
	return report, nil
}

func (r IntegrityReport) Healthy() bool {
	return len(r.Issues) == 0
}

func (r IntegrityReport) EntityCount() int {
	return r.Records + r.Events + r.Workflows + r.Attachments + r.LedgerEntries
}

func (s *Store) EnsureRecordStatus(id string, status model.RecordStatus) error {
	record, err := s.GetRecord(id)
	if err != nil {
		return err
	}
	if !model.ValidStatus(status) {
		return fmt.Errorf("invalid status %s", status)
	}
	record.Status = status
	return s.ReplaceRecord(record)
}
