package report

import (
	"strings"
	"testing"

	"inventorychain/internal/model"
)

func TestSummaryAndCSV(t *testing.T) {
	record := model.Record{ID: "r", Warehouse: "north", Cycle: "c", Status: model.StatusApproved, Lines: []model.DiscrepancyLine{{SKU: "sku", Expected: 1, Observed: 3, Adjustment: 2, Verified: true}}}
	summary := BuildSummary([]model.Record{record})
	if summary.LineCount != 1 || summary.PositiveTotal != 2 || NeedsAttention(summary) {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if !strings.Contains(CSV([]model.Record{record}), "sku") {
		t.Fatal("csv did not include sku")
	}
}
