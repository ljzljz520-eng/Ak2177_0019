package api

import (
	"net/http"
	"strings"

	"inventorychain/internal/service"
)

func CommandHandler(svc *service.Service) http.Handler {
	server := NewServer(svc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/records") {
			server.Handler().ServeHTTP(w, r)
			return
		}
		server.health(w, r)
	})
}
