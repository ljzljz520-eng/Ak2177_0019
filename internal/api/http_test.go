package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"inventorychain/internal/service"
	"inventorychain/internal/store"
)

func TestHTTPHealthAndRecords(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	handler := NewServer(service.New(db, service.FixedClock{Value: "fixed"})).Handler()
	health := httptest.NewRecorder()
	handler.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK || !strings.Contains(health.Body.String(), "ok") {
		t.Fatalf("health response: %d %s", health.Code, health.Body.String())
	}
	list := httptest.NewRecorder()
	handler.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/records?warehouse=north", nil))
	if list.Code != http.StatusOK {
		t.Fatalf("list response: %d", list.Code)
	}
}
