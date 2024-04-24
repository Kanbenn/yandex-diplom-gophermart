package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Kanbenn/gophermart/internal/handler"
)

func New(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(rBodyCloserMiddleware)
	r.Use(gzipMiddleware)

	r.Mount("/debug", chimw.Profiler())

	r.Group(func(r chi.Router) {
		r.Use(requireJsnMiddleware)
		r.Post("/api/user/register", h.RegisterNewUser)
		r.Post("/api/user/login", h.LoginUser)

	})
	r.Group(func(r chi.Router) {
		r.Use(requireAuthMiddleware)
		r.Post("/api/user/orders", h.PostNewOrder)
		r.Get("/api/user/orders", h.GetUserOrders)
		r.Get("/api/user/balance", h.GetUserBalance)
		r.Get("/api/user/withdrawals", h.GetUserWithdrawHistory)
	})
	r.With(requireAuthMiddleware, requireJsnMiddleware).
		Post("/api/user/balance/withdraw", h.PostNewOrderWithBonus)
	return r
}
