package review

import (
	"fmt"

	"inventorychain/internal/model"
)

type Decision struct {
	Allowed bool
	Score   int
	Reasons []string
}

func Evaluate(record model.Record) Decision {
	decision := Decision{Allowed: true, Score: 100, Reasons: []string{}}
	if err := record.Validate(); err != nil {
		decision.Allowed = false
		decision.Score -= 50
		decision.Reasons = append(decision.Reasons, err.Error())
	}
	if record.HasUnverifiedLines() {
		decision.Allowed = false
		decision.Score -= 20
		decision.Reasons = append(decision.Reasons, "unverified discrepancy lines")
	}
	if len(record.Slots) > 1 && !sequenceAligned(record.Slots) {
		decision.Allowed = false
		decision.Score -= 15
		decision.Reasons = append(decision.Reasons, "slot sequence is not aligned")
	}
	if record.DifferenceTotal() == 0 {
		decision.Score += 5
	}
	return decision
}

func sequenceAligned(slots []model.TimeSlot) bool {
	for i, slot := range slots {
		if slot.Sequence != i+1 {
			return false
		}
	}
	return true
}

func RequireApproval(record model.Record) error {
	decision := Evaluate(record)
	if !decision.Allowed {
		return fmt.Errorf("review rejected: %v", decision.Reasons)
	}
	return nil
}

func NormalizeLines(record *model.Record) {
	for i := range record.Lines {
		record.Lines[i].Reason = model.CleanReason(record.Lines[i].Reason)
	}
	record.Recalculate()
}
