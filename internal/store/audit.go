package store

import (
	"sort"

	"inventorychain/internal/model"
)

func (s *Store) SaveEvent(event model.AuditEvent) error {
	return s.put(EventsBucket, event.ID, event)
}

func (s *Store) ListEvents(recordID string) ([]model.AuditEvent, error) {
	events := make([]model.AuditEvent, 0)
	err := s.list(EventsBucket, func(data []byte) error {
		var event model.AuditEvent
		if err := unmarshal(data, &event); err != nil {
			return err
		}
		if event.RecordID == recordID {
			events = append(events, event)
		}
		return nil
	})
	sort.Slice(events, func(i, j int) bool { return events[i].Sequence < events[j].Sequence })
	return events, err
}
