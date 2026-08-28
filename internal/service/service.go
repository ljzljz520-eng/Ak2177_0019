package service

import (
	"fmt"

	"inventorychain/internal/store"
)

type Clock interface {
	Now() string
}

type FixedClock struct {
	Value string
}

func (c FixedClock) Now() string {
	return c.Value
}

type Service struct {
	store *store.Store
	clock Clock
}

func New(st *store.Store, clock Clock) *Service {
	if clock == nil {
		clock = FixedClock{Value: "2026-01-01T00:00:00Z"}
	}
	return &Service{store: st, clock: clock}
}

func (s *Service) Store() *store.Store {
	return s.store
}

func (s *Service) requireStore() error {
	if s == nil || s.store == nil {
		return fmt.Errorf("service store is not configured")
	}
	return nil
}
