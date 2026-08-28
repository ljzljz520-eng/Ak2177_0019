package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"inventorychain/internal/model"
	"inventorychain/internal/service"
)

type Server struct {
	service *service.Service
}

func NewServer(svc *service.Service) *Server {
	return &Server{service: svc}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.record)
	mux.HandleFunc("/records/ledger/", s.ledger)
	mux.HandleFunc("/records/stream/", s.ledgerStream)
	mux.HandleFunc("/imports", s.importRows)
	mux.HandleFunc("/exports", s.exportRows)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var record model.Record
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		created, err := s.service.Register(record, actor(r))
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
		return
	}
	if r.Method == http.MethodGet {
		result, err := s.service.Search(model.SearchFilter{Warehouse: r.URL.Query().Get("warehouse"), Cycle: r.URL.Query().Get("cycle")}, 1, 100)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		writeError(w, http.StatusBadRequest, errMissingID{})
		return
	}
	if r.Method == http.MethodGet {
		record, err := s.service.Get(id)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, record)
		return
	}
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/submit") {
		record, err := s.service.Submit(strings.TrimSuffix(id, "/submit"), actor(r))
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, record)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) ledger(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/records/ledger/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		writeError(w, http.StatusBadRequest, errMissingID{})
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	reconciliation, err := s.service.Reconcile(parts[0])
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, reconciliation)
}

type errMissingID struct{}

func (errMissingID) Error() string { return "record id is required" }

func actor(r *http.Request) string {
	if value := r.Header.Get("X-Actor"); value != "" {
		return value
	}
	return "api"
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
