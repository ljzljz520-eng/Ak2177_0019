package model

import (
	"sort"
	"strings"
)

type LedgerEntry struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	Direction string `json:"direction"`
	Reason    string `json:"reason"`
	Actor     string `json:"actor"`
	Sequence  int    `json:"sequence"`
	Timestamp string `json:"timestamp"`
}

type Reconciliation struct {
	RecordID      string
	ExpectedTotal int
	ObservedTotal int
	NetChange     int
	Balanced      bool
	Entries       []LedgerEntry
}

func NewLedgerEntry(recordID, sku, direction, actor, timestamp string, quantity, sequence int, reason string) LedgerEntry {
	cleanDirection := strings.ToLower(strings.TrimSpace(direction))
	if cleanDirection != "in" && cleanDirection != "out" {
		cleanDirection = "adjustment"
	}
	return LedgerEntry{ID: EventID(recordID, sequence), RecordID: recordID, SKU: sku, Quantity: quantity, Direction: cleanDirection, Reason: CleanReason(reason), Actor: actor, Sequence: sequence, Timestamp: timestamp}
}

func (e LedgerEntry) SignedQuantity() int {
	if e.Direction == "out" {
		return -e.Quantity
	}
	return e.Quantity
}

func (e LedgerEntry) Valid() bool {
	return e.RecordID != "" && e.SKU != "" && e.Quantity >= 0 && e.Sequence > 0 && e.Direction != ""
}

func BuildReconciliation(record Record, entries []LedgerEntry) Reconciliation {
	reconciliation := Reconciliation{RecordID: record.ID, Entries: append([]LedgerEntry(nil), entries...)}
	for _, line := range record.Lines {
		reconciliation.ExpectedTotal += line.Expected
		reconciliation.ObservedTotal += line.Observed
	}
	for _, entry := range entries {
		reconciliation.NetChange += entry.SignedQuantity()
	}
	reconciliation.Balanced = reconciliation.ObservedTotal-reconciliation.ExpectedTotal == reconciliation.NetChange
	return reconciliation
}

func SortLedger(entries []LedgerEntry) []LedgerEntry {
	ordered := append([]LedgerEntry(nil), entries...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Sequence == ordered[j].Sequence {
			return ordered[i].SKU < ordered[j].SKU
		}
		return ordered[i].Sequence < ordered[j].Sequence
	})
	return ordered
}

func SlotOrder(slots []TimeSlot) []string {
	ordered := append([]TimeSlot(nil), slots...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Sequence < ordered[j].Sequence })
	names := make([]string, 0, len(ordered))
	for _, slot := range ordered {
		names = append(names, slot.Name)
	}
	return names
}

func CloneRecord(record Record) Record {
	clone := record
	clone.Slots = append([]TimeSlot(nil), record.Slots...)
	clone.Lines = append([]DiscrepancyLine(nil), record.Lines...)
	return clone
}
