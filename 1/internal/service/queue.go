package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wb-go/wbf/zlog"
)

func (s *Service) PublishNotifications(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			notifications, err := s.repo.CheckNotifications()
			if err != nil {
				return err
			}
			for _, n := range notifications {
				marshalledN, err := json.Marshal(n)
				if err != nil {
					return err
				}
				if err := s.Rabbit.Publisher.Publish(ctx, marshalledN, "notifications"); err != nil {
					return err
				}
				if err := s.repo.UpdateNotificationStatus(n.ID, "queued"); err != nil {
					zlog.Logger.Error().Msgf("уведомление отправлено в очередь, но статус не обновился для notification_id=%d: %v", n.ID, err)
				}
			}
		}
	}
}
