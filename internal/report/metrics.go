package report

import (
	"sort"

	"inventorychain/internal/model"
)

type WarehouseMetric struct {
	Warehouse string
	Records   int
	Lines     int
	Variance  int
	Published int
}

type CycleMetric struct {
	Cycle    string
	Records  int
	Approved int
	Archived int
	Open     int
}

func WarehouseMetrics(records []model.Record) []WarehouseMetric {
	byWarehouse := make(map[string]WarehouseMetric)
	for _, record := range records {
		metric := byWarehouse[record.Warehouse]
		metric.Warehouse = record.Warehouse
		metric.Records++
		if record.Published {
			metric.Published++
		}
		for _, line := range record.Lines {
			metric.Lines++
			metric.Variance += line.Adjustment
		}
		byWarehouse[record.Warehouse] = metric
	}
	metrics := make([]WarehouseMetric, 0, len(byWarehouse))
	for _, metric := range byWarehouse {
		metrics = append(metrics, metric)
	}
	sort.Slice(metrics, func(i, j int) bool { return metrics[i].Warehouse < metrics[j].Warehouse })
	return metrics
}

func CycleMetrics(records []model.Record) []CycleMetric {
	byCycle := make(map[string]CycleMetric)
	for _, record := range records {
		metric := byCycle[record.Cycle]
		metric.Cycle = record.Cycle
		metric.Records++
		switch record.Status {
		case model.StatusApproved:
			metric.Approved++
		case model.StatusArchived:
			metric.Archived++
		default:
			metric.Open++
		}
		byCycle[record.Cycle] = metric
	}
	metrics := make([]CycleMetric, 0, len(byCycle))
	for _, metric := range byCycle {
		metrics = append(metrics, metric)
	}
	sort.Slice(metrics, func(i, j int) bool { return metrics[i].Cycle < metrics[j].Cycle })
	return metrics
}

func Matrix(records []model.Record) map[string]map[model.RecordStatus]int {
	matrix := make(map[string]map[model.RecordStatus]int)
	for _, record := range records {
		if matrix[record.Warehouse] == nil {
			matrix[record.Warehouse] = make(map[model.RecordStatus]int)
		}
		matrix[record.Warehouse][record.Status]++
	}
	return matrix
}

func CompletionRate(records []model.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	complete := 0
	for _, record := range records {
		if record.Status == model.StatusArchived || record.Published {
			complete++
		}
	}
	return float64(complete) / float64(len(records))
}
