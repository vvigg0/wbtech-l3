package service

func (s *Service) DeleteNotification(id int) error {
	return s.repo.DeleteNotification(id)
}
