package worker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/go-resty/resty/v2"
)

type Worker struct {
	newOrderCh chan models.Order
	cfg        config.Config
	s          storer
}

type storer interface {
	UpdateOrderStatusAndUserBalance(order models.AccrualResponse)
}

func New(cfg config.Config, s storer) *Worker {
	ch := make(chan models.Order, 1)
	return &Worker{ch, cfg, s}
}

func (worker *Worker) NotifyWorker(order models.Order) {
	go func() {
		worker.newOrderCh <- order
	}()
}

func (worker *Worker) LaunchWorkerAccrual(ctx context.Context) {
	log.Println("launching accrual worker")
	httpc := resty.New().SetBaseURL(worker.cfg.AccrualLink)

	for {
		select {
		case <-ctx.Done():
			log.Println("worker: recieved ctx.Done() signal, exiting the goroutine")
			return
		case order := <-worker.newOrderCh:
			worker.askAccrual(httpc, order)

		}
	}
}

func (worker *Worker) askAccrual(httpc *resty.Client, in models.Order) {
	var out models.AccrualResponse
	out.User = in.User
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	accrualPath := "/api/orders/" + in.Number
	req := httpc.R().
		SetContext(ctx).
		SetResult(&out)
	resp, err := req.Get(accrualPath)
	cancel()

	log.Println("worker: response from Accrual: \n",
		worker.cfg.AccrualLink+accrualPath,
		in, out, resp.StatusCode(), err)

	if resp.StatusCode() != http.StatusOK {
		return
	}
	worker.s.UpdateOrderStatusAndUserBalance(out)
}
