package model

import "testing"

func TestRecordValidationAndRecalculate(t *testing.T) {
	record := Record{ID: "r1", Warehouse: "north", Cycle: "c1", Slots: []TimeSlot{{Name: "morning", Sequence: 1}}, Lines: []DiscrepancyLine{{SKU: "sku-1", Expected: 4, Observed: 6, Reason: "count variance"}}}
	record.Recalculate()
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
	if record.DifferenceTotal() != 2 || !record.Lines[0].Verified {
		t.Fatalf("unexpected calculated line: %+v", record.Lines[0])
	}
}
