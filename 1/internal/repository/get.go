package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/model"
)

func (r *Repository) GetStatus(id int) (*dto.NotificationStatusDTO, error) {
	query := `SELECT id,status from notifications WHERE id=$1`

	var status dto.NotificationStatusDTO
	err := r.Master.QueryRow(
		query,
		id).Scan(&status.ID, &status.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoNotification
		}
		return nil, fmt.Errorf("не удалось получить статус уведомления из БД: %w", err)
	}

	return &status, nil
}

func (r *Repository) CheckNotifications() ([]model.Notification, error) {
	now := time.Now()
	getQuery := `SELECT *
			FROM notifications 
			WHERE status='active' AND send_at<=$1`

	rows, err := r.Master.Query(getQuery, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifs := []model.Notification{}
	for rows.Next() {
		var notif model.Notification
		if err := rows.Scan(&notif.ID, &notif.Text, &notif.TgID, &notif.Status, &notif.SendAt); err != nil {
			return nil, err
		}
		notifs = append(notifs, notif)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifs, nil
}
