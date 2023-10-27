package storage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/go-resty/resty/v2"
)

func (pg *Pg) LaunchWorkerAccrual() {
	log.Println("launching accrual worker")

	for {
		orders, err := pg.selectOrdersUnfinished()
		if err != nil {
			log.Println("Pg.selectOrdersUnfinished error:", orders, err)
		}
		if len(orders) < 1 {
			continue
		}
		for _, order := range orders {
			result, err := pg.askAccrualForOrderUpdates(order)
			if err != nil {
				log.Println("Pg.Worker request-Accrual error:", order, err)
				time.Sleep(10 * time.Second)
				continue
			}
			log.Println("Pg.Worker updating order status:", order, result)
			if result.Status != order.Status {
				pg.updateOrderStatusAndUserBalance(result)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (pg *Pg) askAccrualForOrderUpdates(order models.AccrualResponse) (models.AccrualResponse, error) {
	httpc := resty.New().SetBaseURL(pg.Cfg.AccrualLink)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	accrualPath := "/api/orders/" + order.Number
	log.Println("making a request to Accrual:", pg.Cfg.AccrualLink, order)
	req := httpc.R().
		SetContext(ctx).
		SetResult(&order)
	resp, err := req.Get(accrualPath)
	cancel()
	log.Println(
		"Pg.Worker response from accrual:",
		accrualPath,
		order,
		resp.StatusCode(),
		err)
	if resp.StatusCode() != http.StatusOK {
		return order, err
	}
	return order, nil
}

func (pg *Pg) updateOrderStatusAndUserBalance(order models.AccrualResponse) {
	tx, err := pg.Sqlx.Beginx()
	if err != nil {
		log.Println("Pg.Worker error updating order in db:", err)
	}
	defer tx.Rollback()

	qo := `UPDATE orders SET status=:status, bonus=:bonus WHERE number=:number`
	if _, err := tx.NamedExec(qo, order); err != nil {
		log.Println("Pg.Worker error updating order in db:", err)
	}
	if order.Status == "PROCESSED" {
		qb := `UPDATE users SET balance= users.balance + :bonus WHERE id=:user_id`
		if _, err = tx.NamedExec(qb, order); err != nil {
			log.Println("Pg.Worker error updating order in db:", err)
		}
	}
	tx.Commit()
	log.Println("Pg.Worker order updated in db:", order)
}

func (pg *Pg) selectOrdersUnfinished() (orders []models.AccrualResponse, err error) {
	q := "SELECT number, user_id, status FROM orders WHERE status NOT IN ('%s')"
	finishedStatuses := strings.Join(pg.Cfg.FinishedOrderStatuses, "','")
	q = fmt.Sprintf(q, finishedStatuses)
	err = pg.Sqlx.Select(&orders, q)
	return orders, err
}
