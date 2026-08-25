package service

import (
	"fmt"
	"sort"

	"inventorychain/internal/model"
)

type PolicySummary struct {
	PolicyName string
	Passed     int
	Failed     int
	Scores     []int
	Violations map[string][]string
}

func (s *Service) EvaluatePolicy(recordIDs []string) (PolicySummary, error) {
	policy := model.DefaultPolicy()
	if err := model.ValidatePolicy(policy); err != nil {
		return PolicySummary{}, err
	}
	summary := PolicySummary{PolicyName: policy.Name, Scores: []int{}, Violations: make(map[string][]string)}
	for _, id := range recordIDs {
		record, err := s.Get(id)
		if err != nil {
			return PolicySummary{}, err
		}
		result := policy.Evaluate(record)
		summary.Scores = append(summary.Scores, result.Score)
		if result.Passed {
			summary.Passed++
		} else {
			summary.Failed++
			summary.Violations[id] = append([]string(nil), result.Violations...)
		}
	}
	sort.Ints(summary.Scores)
	return summary, nil
}

func (s *Service) RequirePolicy(recordID string) error {
	summary, err := s.EvaluatePolicy([]string{recordID})
	if err != nil {
		return err
	}
	if summary.Failed > 0 {
		return fmt.Errorf("record %s did not satisfy review policy", recordID)
	}
	return nil
}

func (s *Service) PolicyName() string {
	return model.DefaultPolicy().Name
}

func (s *Service) PolicyScore(recordID string) (int, []string, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return 0, nil, err
	}
	result := model.DefaultPolicy().Evaluate(record)
	return result.Score, append([]string(nil), result.Violations...), nil
}

func (s *Service) PolicyReady(recordID string) (bool, error) {
	score, violations, err := s.PolicyScore(recordID)
	if err != nil {
		return false, err
	}
	return score >= 75 && len(violations) == 0, nil
}

func PolicyThreshold() int {
	return 75
}

func (p PolicySummary) AverageScore() int {
	if len(p.Scores) == 0 {
		return 0
	}
	total := 0
	for _, score := range p.Scores {
		total += score
	}
	return total / len(p.Scores)
}

func (p PolicySummary) PassedAll() bool {
	return p.Total() > 0 && p.Failed == 0
}

func (p PolicySummary) Total() int {
	return p.Passed + p.Failed
}

func (p PolicySummary) HasViolations(recordID string) bool {
	return len(p.Violations[recordID]) > 0
}

func (p PolicySummary) Label() string {
	if p.PassedAll() {
		return "all records ready"
	}
	return "review required"
}

func (p PolicySummary) ScoreRange() (int, int) {
	if len(p.Scores) == 0 {
		return 0, 0
	}
	return p.Scores[0], p.Scores[len(p.Scores)-1]
}

func (p PolicySummary) FailureRate() float64 {
	if p.Total() == 0 {
		return 0
	}
	return float64(p.Failed) / float64(p.Total())
}

func (p PolicySummary) ReadyCount() int {
	return p.Passed
}

func (p PolicySummary) IsEmpty() bool { return p.Total() == 0 }

func (p PolicySummary) ReadyRatio() float64 { return float64(p.Passed) / float64(maxInt(1, p.Total())) }

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}
