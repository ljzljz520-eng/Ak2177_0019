package flow015

import (
	"reflect"
	"testing"

	"inventorychain/internal/model"
)

func Test2177BusinessRegression(t *testing.T) {
	processor := NewProcessor()
	input := []model.TimeSlot{{Name: "morning", Sequence: 1, Owner: "alice"}, {Name: "afternoon", Sequence: 2, Owner: "bob"}}
	got := processor.Allocate(input).Values()
	if !reflect.DeepEqual(got, input) {
		t.Fatalf("slot order changed: got %v want %v", got, input)
	}
}

func TestSingleSlotAllocation(t *testing.T) {
	processor := NewProcessor()
	input := []model.TimeSlot{{Name: "night", Sequence: 1}}
	got := processor.Allocate(input).Names()
	if !reflect.DeepEqual(got, []string{"night"}) {
		t.Fatalf("unexpected names: %v", got)
	}
}
