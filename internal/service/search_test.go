package service

import (
	"path/filepath"
	"testing"

	"inventorychain/internal/model"
	"inventorychain/internal/store"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "search.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db, FixedClock{Value: "fixed"})
	record, err := svc.Register(model.Record{Warehouse: "south", Cycle: "cycle-2", Slots: []model.TimeSlot{{Name: "day", Sequence: 1}}, Lines: []model.DiscrepancyLine{{SKU: "old", Expected: 1, Observed: 1}}}, "user")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Submit(record.ID, "user"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Approve(record.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	result, err := svc.Search(model.SearchFilter{Warehouse: "south"}, 1, 10)
	if err != nil || result.Total != 1 {
		t.Fatalf("search result: %+v %v", result, err)
	}
	if _, err = svc.Publish(record.ID, "publisher"); err != nil {
		t.Fatal(err)
	}
	got, err := svc.Get(record.ID)
	if err != nil || !got.Published {
		t.Fatalf("publish result: %+v %v", got, err)
	}
}
