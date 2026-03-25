package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/model"
	"github.com/wb-go/wbf/rabbitmq"
)

type repoMock struct {
	createFn func(text string, TgID int64, date time.Time) (int, error)
	checkFn  func() ([]model.Notification, error)
	cancelFn func(id int) error
	updateFn func(id int, newStatus string) error
	getFn    func(id int) (*dto.NotificationStatusDTO, error)
}

func (m *repoMock) Create(text string, TgID int64, date time.Time) (int, error) {
	return m.createFn(text, TgID, date)
}

func (m *repoMock) CheckNotifications() ([]model.Notification, error) {
	return m.checkFn()
}

func (m *repoMock) CancelNotification(id int) error {
	return m.cancelFn(id)
}

func (m *repoMock) UpdateNotificationStatus(id int, newStatus string) error {
	return m.updateFn(id, newStatus)
}

func (m *repoMock) GetStatus(id int) (*dto.NotificationStatusDTO, error) {
	return m.getFn(id)
}

type publisherMock struct {
	publishFn func(ctx context.Context, body []byte, key string) error

	calls  int
	bodies [][]byte
	keys   []string
}

func (m *publisherMock) Publish(
	ctx context.Context,
	body []byte,
	key string,
	opts ...rabbitmq.PublishOption,
) error {
	m.calls++
	cp := make([]byte, len(body))
	copy(cp, body)
	m.bodies = append(m.bodies, cp)
	m.keys = append(m.keys, key)
	return m.publishFn(ctx, body, key)
}

func TestService_CreateNotifications(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		req         dto.CreateNotificationsRequest
		createErr   error
		wantErr     bool
		wantErrText string
		wantIDs     []int
		wantCalls   int
	}{
		{
			name: "success",
			req: dto.CreateNotificationsRequest{
				Notifs: []dto.NotificationDTO{
					{Text: "one", TgID: 1, SendAt: dto.CustomTime{Time: now}},
					{Text: "two", TgID: 2, SendAt: dto.CustomTime{Time: now}},
				},
			},
			wantErr:   false,
			wantIDs:   []int{10, 11},
			wantCalls: 2,
		},
		{
			name: "invalid notifications return validation error",
			req: dto.CreateNotificationsRequest{
				Notifs: []dto.NotificationDTO{
					{Text: "", TgID: 1, SendAt: dto.CustomTime{Time: now}},
					{Text: "ok", TgID: 2, SendAt: dto.CustomTime{Time: now}},
					{Text: "bad tg", TgID: 0, SendAt: dto.CustomTime{Time: now}},
				},
			},
			wantErr:     true,
			wantErrText: "заполнены не все поля",
			wantCalls:   1,
			wantIDs:     []int{10},
		},
		{
			name: "empty request",
			req: dto.CreateNotificationsRequest{
				Notifs: nil,
			},
			wantErr:     true,
			wantErrText: "нужно хотя бы одно уведомление",
			wantIDs:     nil,
			wantCalls:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createCalls := 0

			repo := &repoMock{
				createFn: func(text string, TgID int64, date time.Time) (int, error) {
					createCalls++
					if tt.createErr != nil {
						return 0, tt.createErr
					}
					return 9 + createCalls, nil
				},
				checkFn:  func() ([]model.Notification, error) { return nil, nil },
				cancelFn: func(id int) error { return nil },
				updateFn: func(id int, newStatus string) error { return nil },
				getFn:    func(id int) (*dto.NotificationStatusDTO, error) { return nil, nil },
			}

			pub := &publisherMock{
				publishFn: func(ctx context.Context, body []byte, key string) error { return nil },
			}

			s := New(repo, pub)

			gotIDs, err := s.CreateNotifications(context.Background(), tt.req)

			if tt.wantErr && err == nil {
				t.Fatal("ожидалась ошибка, получили nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("не ожидалась ошибка, получили: %v", err)
			}

			if tt.wantErrText != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("ожидали ошибку с подстрокой %q, получили: %v", tt.wantErrText, err)
				}
			}

			if len(gotIDs) != len(tt.wantIDs) {
				t.Fatalf("ожидали %d id, получили %d", len(tt.wantIDs), len(gotIDs))
			}

			for i := range gotIDs {
				if gotIDs[i] != tt.wantIDs[i] {
					t.Fatalf("ожидали id[%d]=%d, получили %d", i, tt.wantIDs[i], gotIDs[i])
				}
			}

			if createCalls != tt.wantCalls {
				t.Fatalf("ожидали %d вызовов Create, получили %d", tt.wantCalls, createCalls)
			}
		})
	}
}

func TestService_GetNotificationStatus(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		repoResp *dto.NotificationStatusDTO
		repoErr  error
		wantErr  bool
		wantNil  bool
	}{
		{
			name: "success",
			id:   10,
			repoResp: &dto.NotificationStatusDTO{
				ID:     10,
				Status: "queued",
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name:    "repo error",
			id:      20,
			repoErr: errors.New("not found"),
			wantErr: true,
			wantNil: true,
		},
		{
			name:     "nil response without error",
			id:       30,
			repoResp: nil,
			wantErr:  false,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoMock{
				createFn: func(text string, telegramID int64, date time.Time) (int, error) { return 0, nil },
				checkFn:  func() ([]model.Notification, error) { return nil, nil },
				cancelFn: func(id int) error { return nil },
				updateFn: func(id int, newStatus string) error { return nil },
				getFn: func(id int) (*dto.NotificationStatusDTO, error) {
					return tt.repoResp, tt.repoErr
				},
			}

			pub := &publisherMock{
				publishFn: func(ctx context.Context, body []byte, key string) error { return nil },
			}

			s := New(repo, pub)

			got, err := s.GetNotificationStatus(tt.id)

			if tt.wantErr && err == nil {
				t.Fatal("ожидалась ошибка, получили nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("не ожидалась ошибка, получили: %v", err)
			}

			if tt.wantNil && got != nil {
				t.Fatalf("ожидали nil, получили %+v", got)
			}

			if !tt.wantNil && got == nil {
				t.Fatal("ожидали dto, получили nil")
			}

			if !tt.wantNil && got.ID != tt.repoResp.ID {
				t.Fatalf("ожидали id=%d, получили %d", tt.repoResp.ID, got.ID)
			}
		})
	}
}

func TestService_CancelNotification(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		cancelErr error
		wantErr   bool
		wantCalls int
	}{
		{
			name:      "success",
			id:        1,
			wantErr:   false,
			wantCalls: 1,
		},
		{
			name:      "repo error",
			id:        2,
			cancelErr: errors.New("not found"),
			wantErr:   true,
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancelCalls := 0

			repo := &repoMock{
				createFn: func(text string, telegramID int64, date time.Time) (int, error) { return 0, nil },
				checkFn:  func() ([]model.Notification, error) { return nil, nil },
				cancelFn: func(id int) error {
					cancelCalls++
					return tt.cancelErr
				},
				updateFn: func(id int, newStatus string) error { return nil },
				getFn:    func(id int) (*dto.NotificationStatusDTO, error) { return nil, nil },
			}

			pub := &publisherMock{
				publishFn: func(ctx context.Context, body []byte, key string) error { return nil },
			}

			s := New(repo, pub)

			err := s.CancelNotification(tt.id)

			if tt.wantErr && err == nil {
				t.Fatal("ожидалась ошибка, получили nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("не ожидалась ошибка, получили: %v", err)
			}

			if cancelCalls != tt.wantCalls {
				t.Fatalf("ожидали %d вызовов CancelNotification, получили %d", tt.wantCalls, cancelCalls)
			}
		})
	}
}
