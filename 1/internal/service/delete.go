package service

func (s *Service) UpdateNotificationStatus(id int, newStatus string) error {
	return s.repo.UpdateNotificationStatus(id, newStatus)
}
