package worker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/go-resty/resty/v2"
)

func (worker *Worker) LaunchWorkerAccrual(ctx context.Context) {
	log.Println("launching accrual worker goroutine")

	// ticker := time.NewTicker(100 * time.Millisecond)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker: recieved ctx.Done signal, stopping the goroutine")
			return
		default:
			<-ticker.C
			worker.doWork()
		}
	}
}

func (worker *Worker) doWork() {
	orders, err := worker.s.SelectOrdersForAccrual()
	if err != nil {
		log.Println("worker: error at selecting orders to process", err)
		return
	}
	if len(orders) < 1 {
		return
	}
	for _, order := range orders {
		result := worker.askAccrualForOrderUpdates(&order)
		if result.statusCode == http.StatusOK {
			worker.s.UpdateOrderStatusAndUserBalance(order)
		}
		if result.statusCode == http.StatusTooManyRequests {
			log.Println("worker: too many requests to Accrual, wait for", result.delay)
			time.Sleep(result.delay)
		}
	}
}

func (worker *Worker) askAccrualForOrderUpdates(order *models.Accrual) (res result) {
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 900*time.Millisecond)
	accrualPath := "/api/orders/" + order.Number
	httpc := resty.New().SetBaseURL(worker.cfg.AccrualLink)
	req := httpc.R().
		SetContext(ctx).
		SetResult(&order)
	resp, err := req.Get(accrualPath)
	cancel()

	log.Println("worker: response from Accrual: \n",
		worker.cfg.AccrualLink+accrualPath,
		order, resp.StatusCode(), err)

	if resp.StatusCode() == http.StatusNoContent {
		log.Println("worker: order not found in Accrual", order)
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		res.delay = parseSeconds(resp.Header().Get("Retry-After"))
	}
	res.statusCode = resp.StatusCode()
	return res
}

func parseSeconds(seconds string) time.Duration {
	dur, err := time.ParseDuration(seconds + "s")
	if err != nil || dur < 1 {
		log.Println("worker: error at parsing delay from Accrual", err, dur)
		return 60 * time.Second
	}
	return dur
}
