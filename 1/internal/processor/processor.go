package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/model"
)

type TelegramSender interface {
	SendNotification(chatID int64, text string) error
}
type Job struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type Processor struct {
	tg TelegramSender
}

func New(tg TelegramSender) *Processor {
	return &Processor{tg: tg}
}

func (p *Processor) HandleMessage(ctx context.Context, d amqp091.Delivery) error {
	var notif model.Notification
	if err := json.Unmarshal(d.Body, &notif); err != nil {
		return fmt.Errorf("ошибка unmarshall сообщения: %w", err)
	}

	if err := p.tg.SendNotification(notif.TgID, notif.Text); err != nil {
		return fmt.Errorf("tg bot ошибка отправки сообщения: %w", err)
	}
	log.Printf("уведомление отправлено")
	return nil
}
