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
		r.Post("/api/user/login", h.AuthUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(RequireAuthMiddleware)
		r.Post("/api/user/orders", h.PostOrder)
		r.Get("/api/user/orders", h.GetOrders)
		// 	r.Get("/balance", h.GetBalance)
		// 	r.Post("/balance/withdraw", h.PostWithdraw)
		// 	r.Get("/withdrawals", h.GetWithdrawals)
	})

	// POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
	// GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
	// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	// GET /api/user/withdrawals

	return r
}
