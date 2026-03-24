package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/model"
	"github.com/wb-go/wbf/zlog"
)

type NotificationStatusUpdater interface {
	UpdateNotificationStatus(id int, newStatus string) error
}
type TelegramSender interface {
	SendNotification(chatID int64, text string) error
}

type Processor struct {
	tg      TelegramSender
	updater NotificationStatusUpdater
}

func New(tg TelegramSender, updater NotificationStatusUpdater) *Processor {
	return &Processor{tg: tg, updater: updater}
}

func (p *Processor) HandleMessage(ctx context.Context, d amqp091.Delivery) error {
	var notif model.Notification
	if err := json.Unmarshal(d.Body, &notif); err != nil {
		return fmt.Errorf("ошибка unmarshall сообщения: %w", err)
	}

	if err := p.tg.SendNotification(notif.TgID, notif.Text); err != nil {
		return fmt.Errorf("tg bot ошибка отправки сообщения: %w", err)
	}
	if err := p.updater.UpdateNotificationStatus(notif.ID, "sent"); err != nil {
		zlog.Logger.Error().Msgf("уведомление отправлено пользователю, но статус не обновился для notification_id=%d: %v", notif.ID, err)
	}
	return nil
}
