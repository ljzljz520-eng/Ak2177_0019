package model

import (
	"fmt"
	"strings"
)

type PolicyRule struct {
	Name      string
	Threshold int
	Required  bool
	Message   string
}

type ReviewPolicy struct {
	Name          string
	Rules         []PolicyRule
	RequireReason bool
	AllowZeroDiff bool
}

type PolicyResult struct {
	Passed     bool
	Violations []string
	Score      int
}

func DefaultPolicy() ReviewPolicy {
	return ReviewPolicy{Name: "warehouse-cycle-standard", RequireReason: true, AllowZeroDiff: true, Rules: []PolicyRule{{Name: "line-count", Threshold: 1, Required: true, Message: "a discrepancy line is required"}, {Name: "slot-count", Threshold: 1, Required: true, Message: "a counting slot is required"}, {Name: "variance", Threshold: 1000, Required: false, Message: "variance exceeds review threshold"}}}
}

func (p ReviewPolicy) Evaluate(record Record) PolicyResult {
	result := PolicyResult{Passed: true, Score: 100, Violations: []string{}}
	for _, rule := range p.Rules {
		value := p.ruleValue(rule.Name, record)
		if rule.Required && value < rule.Threshold {
			result.Passed = false
			result.Score -= 30
			result.Violations = append(result.Violations, rule.Message)
		}
		if !rule.Required && value > rule.Threshold {
			result.Score -= 10
			result.Violations = append(result.Violations, rule.Message)
		}
	}
	if p.RequireReason {
		for _, line := range record.Lines {
			if line.Adjustment != 0 && strings.TrimSpace(line.Reason) == "" {
				result.Passed = false
				result.Score -= 20
				result.Violations = append(result.Violations, fmt.Sprintf("reason missing for %s", line.SKU))
			}
		}
	}
	if !p.AllowZeroDiff && record.DifferenceTotal() == 0 {
		result.Passed = false
		result.Score -= 10
		result.Violations = append(result.Violations, "zero difference is not allowed")
	}
	return result
}

func (p ReviewPolicy) ruleValue(name string, record Record) int {
	switch name {
	case "line-count":
		return len(record.Lines)
	case "slot-count":
		return len(record.Slots)
	case "variance":
		value := record.DifferenceTotal()
		if value < 0 {
			return -value
		}
		return value
	default:
		return 0
	}
}

func (r PolicyResult) Summary() string {
	if r.Passed {
		return fmt.Sprintf("passed with score %d", r.Score)
	}
	return fmt.Sprintf("failed with score %d: %s", r.Score, strings.Join(r.Violations, "; "))
}

func ValidatePolicy(policy ReviewPolicy) error {
	if strings.TrimSpace(policy.Name) == "" {
		return fmt.Errorf("policy name is required")
	}
	if len(policy.Rules) == 0 {
		return fmt.Errorf("at least one policy rule is required")
	}
	for _, rule := range policy.Rules {
		if strings.TrimSpace(rule.Name) == "" || rule.Threshold < 0 {
			return fmt.Errorf("invalid policy rule")
		}
	}
	return nil
}
