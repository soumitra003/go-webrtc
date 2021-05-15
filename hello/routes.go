package hello

import "github.com/go-chi/chi"

func (h *ModuleHello) RegisterRoutes(router chi.Router) {
	router.Get("/hello/sayhello", h.sayHelloHandler)
}
