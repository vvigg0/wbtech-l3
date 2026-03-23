package repository

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
