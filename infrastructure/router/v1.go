package router

import (
	"net/http"

	"github.com/panuwat39/my-api/internal/shared/response"
)

func registerV1Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/health", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		_ = response.Error(
			w,
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"method not allowed",
		)
		return
	}

	_ = response.Success(
		w,
		http.StatusOK,
		map[string]string{
			"status":  "ok",
			"version": "v1",
		},
	)
}
