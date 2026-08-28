package model

import "testing"

func TestFilterAndPagination(t *testing.T) {
	records := []Record{{ID: "1", Warehouse: "north", Cycle: "a", Status: StatusDraft}, {ID: "2", Warehouse: "south", Cycle: "b", Status: StatusApproved}}
	filter := SearchFilter{Warehouse: "north"}
	if !filter.Matches(records[0]) || filter.Matches(records[1]) {
		t.Fatal("filter mismatch")
	}
	result := Paginate(records, 1, 1)
	if result.Total != 2 || result.Pages != 2 || len(result.Records) != 1 {
		t.Fatalf("unexpected pagination: %+v", result)
	}
}
