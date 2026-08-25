package store

import (
	"path/filepath"
	"testing"

	"inventorychain/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.db")
	record := model.Record{ID: "rec-reopen", Warehouse: "north", Cycle: "cycle-1", Status: model.StatusDraft, Slots: []model.TimeSlot{{Name: "morning", Sequence: 1}}, Lines: []model.DiscrepancyLine{{SKU: "sku", Expected: 1, Observed: 1, Verified: true}}}
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	got, err := second.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Warehouse != record.Warehouse || got.Lines[0].SKU != "sku" {
		t.Fatalf("reopened record mismatch: %+v", got)
	}
}

func TestStoreListsEventsAndAttachments(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.SaveEvent(model.AuditEvent{ID: "e1", RecordID: "r1", Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAttachment(model.Attachment{ID: "a1", RecordID: "r1", Name: "proof"}); err != nil {
		t.Fatal(err)
	}
	events, err := db.ListEvents("r1")
	if err != nil || len(events) != 1 {
		t.Fatalf("events: %v %v", events, err)
	}
	attachments, err := db.ListAttachments("r1")
	if err != nil || len(attachments) != 1 {
		t.Fatalf("attachments: %v %v", attachments, err)
	}
}
