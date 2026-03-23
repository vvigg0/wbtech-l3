package service

import (
	"context"
	"errors"
	"sync"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
	"github.com/wb-go/wbf/zlog"
)

const (
	workers = 3
)

func (s *Service) CreateNotifications(c context.Context, req dto.CreateNotificationsRequest) ([]int, error) {
	var wg sync.WaitGroup
	if len(req.Notifs) == 0 {
		return []int{}, errors.New("нужно хотя бы одно уведомление")
	}

	notifsCh := make(chan dto.NotificationDTO, len(req.Notifs))
	IDsCh := make(chan int, len(req.Notifs))

	for range workers {
		wg.Go(func() { s.worker(c, notifsCh, IDsCh) })
	}

	for _, n := range req.Notifs {
		if n.Text == "" || n.SendAt.IsZero() || n.TgID == 0 {
			continue
		}
		notifsCh <- n
	}

	close(notifsCh)

	go func() {
		wg.Wait()
		close(IDsCh)
	}()

	ids := []int{}
	for id := range IDsCh {
		ids = append(ids, id)
	}

	return ids, nil
}

func (s *Service) worker(ctx context.Context, in chan dto.NotificationDTO, out chan int) {
	for n := range in {
		select {
		case <-ctx.Done():
			return
		default:
			id, err := s.repo.Create(n.Text, n.TgID, n.SendAt.Time)
			if err != nil {
				zlog.Logger.Error().Msgf("ошибка создания уведомления: %v", err)
				continue
			}
			out <- id
		}
	}
}
