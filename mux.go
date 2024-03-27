package entando

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// NewMux return the standard http.ServeMux
// with the default `/api/health` and `/api/metrics` registered.
func NewMux() *http.ServeMux {
	healthStub := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": "up"}`))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthStub)
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
