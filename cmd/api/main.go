package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			log.Printf("write response error: %v", err)
		}
	})

	address := ":8080"

	fmt.Printf("server running on http://localhost%s\n", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
