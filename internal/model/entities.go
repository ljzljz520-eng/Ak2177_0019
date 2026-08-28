package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type RecordStatus string

const (
	StatusDraft    RecordStatus = "draft"
	StatusReview   RecordStatus = "review"
	StatusApproved RecordStatus = "approved"
	StatusArchived RecordStatus = "archived"
	StatusRejected RecordStatus = "rejected"
)

type TimeSlot struct {
	Name     string `json:"name"`
	Sequence int    `json:"sequence"`
	Owner    string `json:"owner"`
}

type DiscrepancyLine struct {
	SKU        string `json:"sku"`
	Expected   int    `json:"expected"`
	Observed   int    `json:"observed"`
	Reason     string `json:"reason"`
	Slot       string `json:"slot"`
	Adjustment int    `json:"adjustment"`
	Verified   bool   `json:"verified"`
}

type Record struct {
	ID        string            `json:"id"`
	Warehouse string            `json:"warehouse"`
	Cycle     string            `json:"cycle"`
	Status    RecordStatus      `json:"status"`
	Slots     []TimeSlot        `json:"slots"`
	Lines     []DiscrepancyLine `json:"lines"`
	CreatedBy string            `json:"created_by"`
	UpdatedBy string            `json:"updated_by"`
	Version   int               `json:"version"`
	Published bool              `json:"published"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Detail    string `json:"detail"`
	Sequence  int    `json:"sequence"`
	Timestamp string `json:"timestamp"`
}

type Workflow struct {
	ID       string   `json:"id"`
	RecordID string   `json:"record_id"`
	Name     string   `json:"name"`
	State    string   `json:"state"`
	Steps    []string `json:"steps"`
	Current  int      `json:"current"`
	Owner    string   `json:"owner"`
}

func SortEvents(events []AuditEvent) []AuditEvent {
	ordered := append([]AuditEvent(nil), events...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Sequence < ordered[j].Sequence })
	return ordered
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Kind     string `json:"kind"`
	Checksum string `json:"checksum"`
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if strings.TrimSpace(r.Warehouse) == "" {
		return errors.New("warehouse is required")
	}
	if strings.TrimSpace(r.Cycle) == "" {
		return errors.New("cycle is required")
	}
	if len(r.Lines) == 0 {
		return errors.New("at least one discrepancy line is required")
	}
	seen := make(map[string]bool)
	for _, slot := range r.Slots {
		if slot.Name == "" || seen[slot.Name] {
			return fmt.Errorf("invalid or duplicate slot %q", slot.Name)
		}
		seen[slot.Name] = true
	}
	for _, line := range r.Lines {
		if line.SKU == "" || line.Expected < 0 || line.Observed < 0 {
			return fmt.Errorf("invalid line %q", line.SKU)
		}
	}
	return nil
}

func (r *Record) Recalculate() {
	for i := range r.Lines {
		r.Lines[i].Adjustment = r.Lines[i].Observed - r.Lines[i].Expected
		r.Lines[i].Verified = r.Lines[i].Adjustment == 0 || r.Lines[i].Reason != ""
	}
}

func (r Record) DifferenceTotal() int {
	total := 0
	for _, line := range r.Lines {
		total += line.Adjustment
	}
	return total
}

func (r Record) HasUnverifiedLines() bool {
	for _, line := range r.Lines {
		if !line.Verified {
			return true
		}
	}
	return false
}

func (w Workflow) Complete() bool {
	return len(w.Steps) > 0 && w.Current >= len(w.Steps)-1 && w.State == "complete"
}
