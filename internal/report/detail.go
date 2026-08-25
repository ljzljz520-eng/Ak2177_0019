package report

import (
	"fmt"
	"sort"
	"strings"

	"inventorychain/internal/model"
)

type Detail struct {
	Record         model.Record
	Events         []model.AuditEvent
	Attachments    []model.Attachment
	Reconciliation model.Reconciliation
}

func BuildDetail(record model.Record, events []model.AuditEvent, attachments []model.Attachment, entries []model.LedgerEntry) Detail {
	return Detail{Record: model.CloneRecord(record), Events: model.SortEvents(events), Attachments: append([]model.Attachment(nil), attachments...), Reconciliation: model.BuildReconciliation(record, entries)}
}

func Timeline(events []model.AuditEvent) []string {
	ordered := model.SortEvents(events)
	lines := make([]string, 0, len(ordered))
	for _, event := range ordered {
		lines = append(lines, fmt.Sprintf("%04d %s %s: %s", event.Sequence, event.Timestamp, event.Action, event.Detail))
	}
	return lines
}

func RenderDetail(detail Detail) string {
	var builder strings.Builder
	builder.WriteString("Record: ")
	builder.WriteString(detail.Record.ID)
	builder.WriteString("\n")
	builder.WriteString("Warehouse: ")
	builder.WriteString(detail.Record.Warehouse)
	builder.WriteString("\n")
	builder.WriteString("Status: ")
	builder.WriteString(StatusLabel(detail.Record.Status))
	builder.WriteString("\n")
	for _, line := range detail.Record.Lines {
		builder.WriteString(fmt.Sprintf("%s %+d verified=%t\n", line.SKU, line.Adjustment, line.Verified))
	}
	for _, item := range Timeline(detail.Events) {
		builder.WriteString(item)
		builder.WriteString("\n")
	}
	return builder.String()
}

func SortAttachments(attachments []model.Attachment) []model.Attachment {
	ordered := append([]model.Attachment(nil), attachments...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Name < ordered[j].Name })
	return ordered
}
