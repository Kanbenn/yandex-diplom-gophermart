package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Kanbenn/gophermart/internal/handler"
	"github.com/Kanbenn/gophermart/internal/storage"
)

func New(h *handler.Handler, s *storage.Pg) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(GzipMiddleware)

	r.Mount("/debug", chimw.Profiler())

	r.Route("/api/user", func(r chi.Router) {
		r.Use(RequireJsnMiddleware)
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.AuthUser)
	})

	// POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
	// GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
	// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	// GET /api/user/withdrawals

	return r
}
