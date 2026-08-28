package model

import (
	"fmt"
	"strings"
)

type ImportSchema struct {
	Required []string
	Optional []string
	Aliases  map[string]string
}

type SchemaIssue struct {
	Field   string
	Message string
	Row     int
}

func DefaultImportSchema() ImportSchema {
	return ImportSchema{Required: []string{"warehouse", "cycle", "sku", "expected", "observed", "slot"}, Optional: []string{"reason"}, Aliases: map[string]string{"location": "warehouse", "item": "sku", "counted": "observed"}}
}

func (s ImportSchema) Normalize(headers []string) []string {
	normalized := make([]string, 0, len(headers))
	for _, header := range headers {
		name := strings.ToLower(strings.TrimSpace(header))
		if alias, ok := s.Aliases[name]; ok {
			name = alias
		}
		normalized = append(normalized, name)
	}
	return normalized
}

func (s ImportSchema) Validate(headers []string) []SchemaIssue {
	actual := make(map[string]bool)
	for _, header := range s.Normalize(headers) {
		actual[header] = true
	}
	issues := make([]SchemaIssue, 0)
	for _, required := range s.Required {
		if !actual[required] {
			issues = append(issues, SchemaIssue{Field: required, Message: "required header missing"})
		}
	}
	return issues
}

func (s SchemaIssue) Error() string {
	if s.Row > 0 {
		return fmt.Sprintf("row %d %s: %s", s.Row, s.Field, s.Message)
	}
	return fmt.Sprintf("%s: %s", s.Field, s.Message)
}

func FormatIssues(issues []SchemaIssue) string {
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, issue.Error())
	}
	return strings.Join(parts, "; ")
}
