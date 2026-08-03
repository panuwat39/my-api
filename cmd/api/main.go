package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/panuwat39/my-api/infrastructure/config"
	"github.com/panuwat39/my-api/infrastructure/router"
)

func main() {
	cfg := config.Load()
	handler := router.New()
	address := ":" + cfg.HTTPPort

	fmt.Printf(
		"server running in %s mode on http://localhost%s\n",
		cfg.AppEnv,
		address,
	)

	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
