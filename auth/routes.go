package auth

import "github.com/go-chi/chi"

// RegisterRoutes registers routes for auth module
func (h *ModuleAuth) RegisterRoutes(router chi.Router) {
	router.Get("/auth/register", h.register)
	router.Get("/auth/google-oauth2-init", h.oAuth2Init)
	router.Get("/login", h.oAuth2Login)
}
