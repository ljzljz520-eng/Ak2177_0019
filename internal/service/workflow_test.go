package service

import (
	"path/filepath"
	"testing"

	"inventorychain/internal/model"
	"inventorychain/internal/store"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db, FixedClock{Value: "2026-01-01T00:00:00Z"})
	record, err := svc.Register(model.Record{Warehouse: "north", Cycle: "cycle-1", Slots: []model.TimeSlot{{Name: "morning", Sequence: 1}}, Lines: []model.DiscrepancyLine{{SKU: "sku", Expected: 4, Observed: 5, Reason: "count variance"}}}, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Submit(record.ID, "alice"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Approve(record.ID, "bob"); err != nil {
		t.Fatal(err)
	}
	archived, err := svc.Archive(record.ID, "bob")
	if err != nil || archived.Status != model.StatusArchived {
		t.Fatalf("archive result: %+v %v", archived, err)
	}
}
