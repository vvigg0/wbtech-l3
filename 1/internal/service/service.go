package service

import (
	"context"
	"time"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/model"
	"github.com/wb-go/wbf/rabbitmq"
)

type notificationRepo interface {
	Create(text string, telegramID int64, date time.Time) (int, error)
	CheckNotifications() ([]model.Notification, error)
	CancelNotification(id int) error
	UpdateNotificationStatus(id int, newStatus string) error
	GetStatus(id int) (*dto.NotificationStatusDTO, error)
}

type publisher interface {
	Publish(ctx context.Context, body []byte, key string, opts ...rabbitmq.PublishOption) error
}
type Service struct {
	repo      notificationRepo
	publisher publisher
}

func New(r notificationRepo, p publisher) *Service {
	return &Service{r, p}
}
