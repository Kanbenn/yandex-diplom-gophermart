package storage

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/go-resty/resty/v2"
)

func (pg *Pg) LaunchWorkerAccrual() {
	log.Println("launching accrual worker")

	for {
		orders, err := pg.SelectOrdersUnprocessed()
		if err != nil {
			log.Println("Pg.selectUnprocessedOrders:", orders, err)
		}
		if len(orders) < 1 {
			continue
		}
		for _, order := range orders {
			result, err := pg.askAccrualForOrderUpdates(order)
			if err != nil {
				log.Println("Pg.Worker request error:", order, err)
				continue
			}
			log.Println("Pg.Worker updating order status:", order, result)
			if err := pg.updateOrderStatusAndUserBalance(result); err != nil {
				log.Println("Pg.Worker error updating order:", order, err)
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

func (pg *Pg) updateOrderStatusAndUserBalance(order models.AccrualResponse) (err error) {
	tx, err := pg.Sqlx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qo := `UPDATE orders SET status=:status, bonus=:bonus WHERE number=:number`
	if _, err := tx.NamedExec(qo, order); err != nil {
		return err
	}
	if order.Status == "PROCESSED" {
		qb := `UPDATE users SET balance= users.balance + :bonus WHERE id=:user_id`
		if _, err = tx.NamedExec(qb, order); err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (pg *Pg) SelectOrdersUnprocessed() (orders []models.AccrualResponse, err error) {
	q := `SELECT number, user_id, status FROM orders
	WHERE status NOT IN ('PROCESSED','INVALID')`
	err = pg.Sqlx.Select(&orders, q)
	return orders, err
}
