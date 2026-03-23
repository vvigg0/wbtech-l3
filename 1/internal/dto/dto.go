package dto

import (
	"strings"
	"time"
)

type NotificationStatusDTO struct {
	ID     int
	Status string
}
type NotificationDTO struct {
	Text   string     `json:"text"`
	TgID   int64      `json:"telegram_ID"`
	SendAt customTime `json:"send_at"`
}

type customTime struct {
	time.Time
}

func (ct *customTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

type CreateNotificationsRequest struct {
	Notifs []NotificationDTO `json:"notifs"`
}
