package http

import "fmt"

// Handler - stores pointer to comment service
type Handler struct { }

// NewHandler - returns a pointer to a Handler
func NewHandler() *Handler{
	return &Handler{}
}

// SetupRoutes - sets up all the routes for the app
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up routes.")
}