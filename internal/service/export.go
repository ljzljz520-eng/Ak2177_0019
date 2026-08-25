package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"inventorychain/internal/model"
)

type ExportEnvelope struct {
	GeneratedAt string         `json:"generated_at"`
	Count       int            `json:"count"`
	Records     []model.Record `json:"records"`
}

type ChangeSet struct {
	RecordID   string
	Before     model.Record
	After      model.Record
	Changed    []string
	VersionGap int
}

func (s *Service) ExportJSON(filter model.SearchFilter) ([]byte, error) {
	result, err := s.Search(filter, 1, 10000)
	if err != nil {
		return nil, err
	}
	envelope := ExportEnvelope{GeneratedAt: s.clock.Now(), Count: result.Total, Records: result.Records}
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode export: %w", err)
	}
	return data, nil
}

func DiffRecords(before, after model.Record) ChangeSet {
	change := ChangeSet{RecordID: after.ID, Before: model.CloneRecord(before), After: model.CloneRecord(after), Changed: []string{}, VersionGap: after.Version - before.Version}
	if before.Status != after.Status {
		change.Changed = append(change.Changed, "status")
	}
	if before.Published != after.Published {
		change.Changed = append(change.Changed, "published")
	}
	if before.Warehouse != after.Warehouse {
		change.Changed = append(change.Changed, "warehouse")
	}
	if before.Cycle != after.Cycle {
		change.Changed = append(change.Changed, "cycle")
	}
	if len(before.Lines) != len(after.Lines) || before.DifferenceTotal() != after.DifferenceTotal() {
		change.Changed = append(change.Changed, "lines")
	}
	if !sameSlots(before.Slots, after.Slots) {
		change.Changed = append(change.Changed, "slots")
	}
	sort.Strings(change.Changed)
	return change
}

func sameSlots(left, right []model.TimeSlot) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func ParseStatus(value string) (model.RecordStatus, error) {
	status := model.RecordStatus(strings.ToLower(strings.TrimSpace(value)))
	if !model.ValidStatus(status) {
		return "", fmt.Errorf("unknown status %q", value)
	}
	return status, nil
}

func (s *Service) RecordsForDashboard() ([]model.Record, error) {
	result, err := s.Search(model.SearchFilter{}, 1, 10000)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(result.Records, func(i, j int) bool {
		if result.Records[i].Status == result.Records[j].Status {
			return result.Records[i].UpdatedAt > result.Records[j].UpdatedAt
		}
		return result.Records[i].Status < result.Records[j].Status
	})
	return result.Records, nil
}
