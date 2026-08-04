package router

import (
	"encoding/json"
	"log"
	"net/http"
)

func registerV1Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/health", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ok",
		"version": "v1",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("write health response error: %v", err)
	}
}
