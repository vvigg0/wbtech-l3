package service

import (
	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
)

func (s *Service) GetNotificationStatus(id int) (*dto.NotificationStatusDTO, error) {
	return s.repo.GetStatus(id)
}
