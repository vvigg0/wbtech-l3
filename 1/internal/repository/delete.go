package repository

import (
	"database/sql"
	"errors"
)

func (r *Repository) UpdateNotificationStatus(id int, newStatus string) error {
	query := "UPDATE notifications SET status=$1 WHERE id=$2"

	res, err := r.Master.Exec(query, newStatus, id)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return ErrNoNotification
	}

	return nil
}

func (r *Repository) CancelNotification(id int) error {
	query := "SELECT status FROM notifications WHERE id=$1"

	var status string
	if err := r.Master.QueryRow(query, id).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoNotification
		}
		return err
	}

	if status == "active" {
		if err := r.UpdateNotificationStatus(id, "cancelled"); err != nil {
			return err
		}
	}

	return nil
}
