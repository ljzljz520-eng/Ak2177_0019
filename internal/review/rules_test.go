package review

import (
	"testing"

	"inventorychain/internal/model"
)

func TestEvaluateAndChecklist(t *testing.T) {
	record := model.Record{ID: "r", Warehouse: "north", Cycle: "c", Slots: []model.TimeSlot{{Name: "morning", Sequence: 1}}, Lines: []model.DiscrepancyLine{{SKU: "sku", Expected: 2, Observed: 3, Reason: "count"}}}
	NormalizeLines(&record)
	decision := Evaluate(record)
	if !decision.Allowed || decision.Score < 100 {
		t.Fatalf("unexpected decision: %+v", decision)
	}
	if !BuildChecklist(record).Ready() {
		t.Fatal("checklist should be ready")
	}
}
