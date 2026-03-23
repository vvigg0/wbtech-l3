package repository

import (
	"time"
)

func (r *Repository) Create(text string, telegramID int64, date time.Time) (int, error) {
	var id int
	err := r.Master.QueryRow(
		`INSERT INTO notifications (text,tg_ID,send_at) VALUES ($1,$2,$3) RETURNING id`,
		text,
		telegramID,
		date).Scan(&id)
	return id, err
}
