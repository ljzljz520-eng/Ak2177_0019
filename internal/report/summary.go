package report

import (
	"fmt"
	"sort"
	"strings"

	"inventorychain/internal/model"
)

type Summary struct {
	Warehouse       string
	Cycle           string
	RecordCount     int
	LineCount       int
	PositiveTotal   int
	NegativeTotal   int
	UnverifiedCount int
	StatusCounts    map[model.RecordStatus]int
}

func BuildSummary(records []model.Record) Summary {
	summary := Summary{StatusCounts: make(map[model.RecordStatus]int)}
	if len(records) > 0 {
		summary.Warehouse = records[0].Warehouse
	}
	for _, record := range records {
		summary.RecordCount++
		summary.StatusCounts[record.Status]++
		if summary.Cycle == "" {
			summary.Cycle = record.Cycle
		}
		for _, line := range record.Lines {
			summary.LineCount++
			if line.Adjustment >= 0 {
				summary.PositiveTotal += line.Adjustment
			} else {
				summary.NegativeTotal += line.Adjustment
			}
			if !line.Verified {
				summary.UnverifiedCount++
			}
		}
	}
	return summary
}

func CSV(records []model.Record) string {
	rows := []string{"record_id,warehouse,cycle,status,sku,adjustment,verified"}
	copyRecords := append([]model.Record(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].ID < copyRecords[j].ID })
	for _, record := range copyRecords {
		for _, line := range record.Lines {
			rows = append(rows, fmt.Sprintf("%s,%s,%s,%s,%s,%d,%t", record.ID, record.Warehouse, record.Cycle, record.Status, line.SKU, line.Adjustment, line.Verified))
		}
	}
	return strings.Join(rows, "\n") + "\n"
}

func StatusLabel(status model.RecordStatus) string {
	switch status {
	case model.StatusDraft:
		return "Draft"
	case model.StatusReview:
		return "In review"
	case model.StatusApproved:
		return "Approved"
	case model.StatusArchived:
		return "Archived"
	default:
		return "Rejected"
	}
}

func NeedsAttention(summary Summary) bool {
	return summary.UnverifiedCount > 0 || summary.NegativeTotal < 0
}
