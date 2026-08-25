package flow015

import (
	"fmt"

	"inventorychain/internal/model"
)

type ReconcileResult struct {
	RecordID string
	Slots    []model.TimeSlot
	Missing  []string
	Extra    []string
	Aligned  bool
}

func ReconcileSlots(expected, actual []model.TimeSlot) ReconcileResult {
	result := ReconcileResult{Slots: append([]model.TimeSlot(nil), actual...), Aligned: true}
	expectedSet := make(map[string]bool, len(expected))
	actualSet := make(map[string]bool, len(actual))
	for index, slot := range expected {
		expectedSet[slot.Name] = true
		if index >= len(actual) || actual[index].Name != slot.Name {
			result.Aligned = false
		}
	}
	for _, slot := range actual {
		actualSet[slot.Name] = true
	}
	for name := range expectedSet {
		if !actualSet[name] {
			result.Missing = append(result.Missing, name)
		}
	}
	for name := range actualSet {
		if !expectedSet[name] {
			result.Extra = append(result.Extra, name)
		}
	}
	return result
}

func EnsureSlotSequence(slots []model.TimeSlot) error {
	if len(slots) == 0 {
		return fmt.Errorf("at least one slot is required")
	}
	seen := make(map[int]bool, len(slots))
	for i, slot := range slots {
		if slot.Sequence != i+1 {
			return fmt.Errorf("slot %s has sequence %d, want %d", slot.Name, slot.Sequence, i+1)
		}
		if seen[slot.Sequence] {
			return fmt.Errorf("duplicate slot sequence %d", slot.Sequence)
		}
		seen[slot.Sequence] = true
	}
	return nil
}

func MergeSlots(primary, secondary []model.TimeSlot) []model.TimeSlot {
	merged := append([]model.TimeSlot(nil), primary...)
	known := make(map[string]bool, len(primary))
	for _, slot := range primary {
		known[slot.Name] = true
	}
	for _, slot := range secondary {
		if !known[slot.Name] {
			merged = append(merged, slot)
			known[slot.Name] = true
		}
	}
	return merged
}

func SlotSummary(slots []model.TimeSlot) map[string]int {
	counts := make(map[string]int)
	for _, slot := range slots {
		counts[slot.Owner]++
	}
	return counts
}
