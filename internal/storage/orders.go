package storage

import (
	"time"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertOrder(o models.OrderInsert) error {
	q := `
	INSERT INTO orders (number, user_id, uploaded_at) VALUES($1, $2, $3)
	ON CONFLICT (number) DO UPDATE SET updated_at = NOW()
	RETURNING user_id, (updated_at > created_at) AS was_conflict`

	o.UploadedAt = time.Now().Format(time.RFC3339)
	result := models.OrderInsertResult{}
	row := pg.Sqlx.QueryRowx(q, o.Number, o.UID, o.UploadedAt)
	result.Err = row.Scan(&result.UID, &result.WasConflict)
	return checkInsertOrderResults(result, o)
}
