package sse

import "github.com/go-chi/chi"

// RegisterRoutes registers routes for auth module
func (h *ModuleSSE) RegisterRoutes(router chi.Router) {
	router.Get("/sse/events", h.subscribeToEvents)
}
