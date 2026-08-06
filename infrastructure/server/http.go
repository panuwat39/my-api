package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type HTTPServer struct {
	app     *fiber.App
	address string
}

func NewHTTPServer(
	port string,
	app *fiber.App,
) *HTTPServer {
	return &HTTPServer{
		app:     app,
		address: ":" + port,
	}
}

func (s *HTTPServer) Start() error {
	if err := s.app.Listen(s.address); err != nil {
		return fmt.Errorf("start Fiber server: %w", err)
	}

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("shutdown Fiber server: %w", err)
	}

	return nil
}
