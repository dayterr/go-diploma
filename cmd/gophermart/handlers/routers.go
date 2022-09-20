package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateRouterWithAsyncHandler(ah AsyncHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(ah.Auth.AuthMiddleware)

	r.Use(middleware.Logger)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", ah.RegisterUser)
		r.Post("/login", ah.LogUser)
		r.Post("/orders", ah.PostOrder)
		r.Get("/orders", ah.LoadOrderList)
		r.Get("/balance", ah.GetBalance)
	})
	return r
}