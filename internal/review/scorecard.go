package review

import (
	"fmt"

	"inventorychain/internal/model"
)

type Scorecard struct {
	Identity     int
	Completeness int
	Evidence     int
	Sequence     int
	Total        int
	Eligible     bool
}

func Score(record model.Record) Scorecard {
	card := Scorecard{}
	if record.ID != "" && record.Warehouse != "" && record.Cycle != "" {
		card.Identity = 25
	}
	if len(record.Lines) > 0 && len(record.Slots) > 0 {
		card.Completeness = 25
	}
	for _, line := range record.Lines {
		if line.Verified || line.Reason != "" {
			card.Evidence += 10
		}
	}
	if card.Evidence > 25 {
		card.Evidence = 25
	}
	if sequenceAligned(record.Slots) {
		card.Sequence = 25
	}
	card.Total = card.Identity + card.Completeness + card.Evidence + card.Sequence
	card.Eligible = card.Total >= 75
	return card
}

func ExplainScore(card Scorecard) []string {
	reasons := make([]string, 0, 4)
	if card.Identity < 25 {
		reasons = append(reasons, "identity fields are incomplete")
	}
	if card.Completeness < 25 {
		reasons = append(reasons, "lines and slots are required")
	}
	if card.Evidence < 25 {
		reasons = append(reasons, "line evidence is incomplete")
	}
	if card.Sequence < 25 {
		reasons = append(reasons, "slot order needs correction")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "record meets review threshold")
	}
	return reasons
}

func ScoreLabel(total int) string {
	switch {
	case total >= 90:
		return "excellent"
	case total >= 75:
		return "ready"
	case total >= 50:
		return "needs evidence"
	default:
		return fmt.Sprintf("blocked-%d", total)
	}
}

func CompareScores(left, right Scorecard) int {
	if left.Total > right.Total {
		return 1
	}
	if left.Total < right.Total {
		return -1
	}
	return 0
}
