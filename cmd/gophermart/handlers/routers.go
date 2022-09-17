package handlers

import "github.com/go-chi/chi/v5"

func CreateRouterWithAsyncHandler(ah AsyncHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", ah.RegisterUser)
		r.Post("/login", ah.LogUser)
	})
	return r
}