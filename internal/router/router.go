package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Kanbenn/gophermart/internal/handler"
)

func New(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(GzipMiddleware)

	r.Mount("/debug", chimw.Profiler())

	r.Group(func(r chi.Router) {
		r.Use(RequireJsnMiddleware)
		r.Post("/api/user/register", h.RegisterUser)
		r.Post("/api/user/login", h.LoginUser)

	})
	r.Group(func(r chi.Router) {
		r.Use(RequireAuthMiddleware)
		r.Post("/api/user/orders", h.PostNewOrder)
		r.Get("/api/user/orders", h.GetUserOrders)
		r.Get("/api/user/balance", h.GetUserBalance)
		r.Get("/api/user/withdrawals", h.GetUserHistory)
	})
	r.Group(func(r chi.Router) {
		r.Use(RequireAuthMiddleware)
		r.Use(RequireJsnMiddleware)
		r.Post("/api/user/balance/withdraw", h.PostNewBonusOrder)
	})
	return r
}
