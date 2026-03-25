package handler

import (
	"context"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
)

type NotificationService interface {
	GetNotificationStatus(id int) (*dto.NotificationStatusDTO, error)
	CreateNotifications(c context.Context, req dto.CreateNotificationsRequest) ([]int, error)
	CancelNotification(id int) error
}
type Handler struct {
	srvc NotificationService
}

func New(s NotificationService) *Handler {
	return &Handler{srvc: s}
}
