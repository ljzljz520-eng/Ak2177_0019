package service

import (
	"path/filepath"
	"testing"

	"inventorychain/internal/store"
)

func TestWorkflowImportReport(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "import.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db, FixedClock{Value: "fixed"})
	result, err := svc.Import([]ImportRow{{Warehouse: "north", Cycle: "batch", SKU: "sku-1", Expected: 2, Observed: 1, Reason: "short", Slot: "morning"}, {Warehouse: "", Cycle: "batch", SKU: "sku-2"}}, "loader")
	if err != nil || len(result.Accepted) != 1 || len(result.Rejected) != 1 {
		t.Fatalf("import result: %+v %v", result, err)
	}
	if result.Summary.LineCount != 1 {
		t.Fatalf("summary: %+v", result.Summary)
	}
}
