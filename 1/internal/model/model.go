package model

import "time"

type Notification struct {
	ID     int       `json:"id"`
	Text   string    `json:"text"`
	TgID   int64     `json:"tg_id"`
	Status string    `json:"status"`
	SendAt time.Time `json:"send_at"`
}
