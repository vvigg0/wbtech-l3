package service

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

func (s *Service) PublishNotifications(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			log.Println("Загружаем уведомления")
			notifications, err := s.repo.CheckNotifications()
			if err != nil {
				return err
			}
			log.Println("нашел ", len(notifications), "уведомлений")
			for _, n := range notifications {
				marshalledN, err := json.Marshal(n)
				if err != nil {
					return err
				}
				if err := s.Rabbit.Publisher.Publish(ctx, marshalledN, "notifications"); err != nil {
					return err
				}
				s.UpdateNotificationStatus(n.ID, "sent")
			}
		}
	}
}
