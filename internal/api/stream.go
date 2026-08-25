package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"

	"inventorychain/internal/model"
	"inventorychain/internal/service"
)

type batchRequest struct {
	Rows  []service.ImportRow `json:"rows"`
	Actor string              `json:"actor"`
}

func (s *Server) importRows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var request batchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.Actor == "" {
		request.Actor = actor(r)
	}
	result, err := s.service.Import(request.Rows, request.Actor)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) exportRows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := s.service.ExportJSON(model.SearchFilter{Warehouse: r.URL.Query().Get("warehouse"), Cycle: r.URL.Query().Get("cycle")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) ledgerStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	recordID := strings.TrimPrefix(r.URL.Path, "/records/")
	if recordID == "" {
		writeError(w, http.StatusBadRequest, errMissingID{})
		return
	}
	rows := bufio.NewScanner(r.Body)
	accepted := 0
	for rows.Scan() {
		parts := strings.Split(rows.Text(), "|")
		if len(parts) < 3 {
			continue
		}
		request := service.AdjustmentRequest{SKU: parts[0], Direction: parts[1], Reason: parts[2]}
		if _, err := s.service.AddAdjustment(recordID, actor(r), request); err == nil {
			accepted++
		}
	}
	writeJSON(w, http.StatusOK, map[string]int{"accepted": accepted})
}
