package flow015

import (
	"fmt"

	"inventorychain/internal/model"
)

type Assignment struct {
	RecordID string
	Slots    []model.TimeSlot
	OrderOK  bool
}

func BuildAssignment(record model.Record) (Assignment, error) {
	if err := record.Validate(); err != nil {
		return Assignment{}, err
	}
	processor := NewProcessor()
	chain := processor.Allocate(record.Slots)
	slots := chain.Values()
	if len(slots) == 0 {
		return Assignment{}, fmt.Errorf("record %s has no slots", record.ID)
	}
	return Assignment{RecordID: record.ID, Slots: slots, OrderOK: ValidateSequence(slots)}, nil
}

func SlotNames(record model.Record) []string {
	processor := NewProcessor()
	return processor.Allocate(record.Slots).Names()
}
