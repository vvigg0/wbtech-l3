package service

func (s *Service) CancelNotification(id int) error {
	return s.repo.CancelNotification(id)
}
