package adminapi

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("/api/dbs", listDBHandler)
	mux.HandleFunc("/api/db/", dbRouter)
}
