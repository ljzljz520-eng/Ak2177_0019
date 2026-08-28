package review

import (
	"fmt"
	"strings"

	"inventorychain/internal/model"
)

type Appeal struct {
	RecordID string
	Actor    string
	Reason   string
	Evidence []string
	Valid    bool
}

func ValidateAppeal(appeal Appeal) error {
	if strings.TrimSpace(appeal.RecordID) == "" {
		return fmt.Errorf("appeal record id is required")
	}
	if strings.TrimSpace(appeal.Actor) == "" {
		return fmt.Errorf("appeal actor is required")
	}
	if strings.TrimSpace(appeal.Reason) == "" {
		return fmt.Errorf("appeal reason is required")
	}
	if len(appeal.Evidence) == 0 {
		return fmt.Errorf("appeal evidence is required")
	}
	return nil
}

func ResolveAppeal(record model.Record, appeal Appeal, approve bool) (model.Record, error) {
	if err := ValidateAppeal(appeal); err != nil {
		return model.Record{}, err
	}
	updated := model.CloneRecord(record)
	if approve {
		updated.Status = model.StatusReview
	} else {
		updated.Status = model.StatusRejected
	}
	updated.UpdatedBy = appeal.Actor
	updated.Version++
	return updated, nil
}

func DecisionSummary(decision Decision) string {
	if decision.Allowed {
		return fmt.Sprintf("approved with score %d", decision.Score)
	}
	return fmt.Sprintf("rejected with score %d: %s", decision.Score, strings.Join(decision.Reasons, "; "))
}
