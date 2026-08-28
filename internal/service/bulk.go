package service

import (
	"fmt"
	"sort"

	"inventorychain/internal/model"
	"inventorychain/internal/review"
)

type BatchResult struct {
	Succeeded []string
	Failed    map[string]string
	Total     int
}

func (s *Service) BulkSubmit(recordIDs []string, actor string) BatchResult {
	result := BatchResult{Succeeded: []string{}, Failed: make(map[string]string), Total: len(recordIDs)}
	for _, id := range recordIDs {
		if _, err := s.Submit(id, actor); err != nil {
			result.Failed[id] = err.Error()
			continue
		}
		result.Succeeded = append(result.Succeeded, id)
	}
	sort.Strings(result.Succeeded)
	return result
}

func (s *Service) BulkApprove(recordIDs []string, actor string) BatchResult {
	result := BatchResult{Succeeded: []string{}, Failed: make(map[string]string), Total: len(recordIDs)}
	for _, id := range recordIDs {
		if _, err := s.Approve(id, actor); err != nil {
			result.Failed[id] = err.Error()
			continue
		}
		result.Succeeded = append(result.Succeeded, id)
	}
	sort.Strings(result.Succeeded)
	return result
}

func (s *Service) ReviewPreview(recordID string) (review.Decision, review.Scorecard, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return review.Decision{}, review.Scorecard{}, err
	}
	decision := review.Evaluate(record)
	card := review.Score(record)
	return decision, card, nil
}

func (s *Service) ValidateBatch(records []model.Record) []string {
	errors := make([]string, 0)
	for _, record := range records {
		if err := record.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", record.ID, err))
			continue
		}
		policy := model.DefaultPolicy()
		if result := policy.Evaluate(record); !result.Passed {
			errors = append(errors, fmt.Sprintf("%s: %s", record.ID, result.Summary()))
		}
	}
	return errors
}

func (s *Service) ArchiveBatch(recordIDs []string, actor string) BatchResult {
	result := BatchResult{Succeeded: []string{}, Failed: make(map[string]string), Total: len(recordIDs)}
	for _, id := range recordIDs {
		if _, err := s.Archive(id, actor); err != nil {
			result.Failed[id] = err.Error()
			continue
		}
		result.Succeeded = append(result.Succeeded, id)
	}
	return result
}
