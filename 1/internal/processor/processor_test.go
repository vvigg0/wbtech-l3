package processor

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/model"
)

type telegramSenderMock struct {
	sendErr error

	called bool
	chatID int64
	text   string
}

func (m *telegramSenderMock) SendNotification(chatID int64, text string) error {
	m.called = true
	m.chatID = chatID
	m.text = text
	return m.sendErr
}

type statusUpdaterMock struct {
	updateErr error

	called bool
	id     int
	status string
}

func (m *statusUpdaterMock) UpdateNotificationStatus(id int, status string) error {
	m.called = true
	m.id = id
	m.status = status
	return m.updateErr
}

func TestProcessor_HandleMessage(t *testing.T) {
	validNotif := model.Notification{
		ID:   1,
		Text: "hello",
		TgID: 123456,
	}

	validBody, err := json.Marshal(validNotif)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	tests := []struct {
		name           string
		body           []byte
		sendErr        error
		updateErr      error
		wantErr        bool
		wantSendCall   bool
		wantUpdateCall bool
	}{
		{
			name:           "invalid json",
			body:           []byte("not-json"),
			wantErr:        true,
			wantSendCall:   false,
			wantUpdateCall: false,
		},
		{
			name:           "telegram send error",
			body:           validBody,
			sendErr:        errors.New("send failed"),
			wantErr:        true,
			wantSendCall:   true,
			wantUpdateCall: false,
		},
		{
			name:           "success",
			body:           validBody,
			wantErr:        false,
			wantSendCall:   true,
			wantUpdateCall: true,
		},
		{
			name:           "update status error ignored",
			body:           validBody,
			updateErr:      errors.New("update failed"),
			wantErr:        false,
			wantSendCall:   true,
			wantUpdateCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgMock := &telegramSenderMock{sendErr: tt.sendErr}
			updaterMock := &statusUpdaterMock{updateErr: tt.updateErr}

			p := New(tgMock, updaterMock)

			err := p.HandleMessage(context.Background(), amqp091.Delivery{
				Body: tt.body,
			})

			if tt.wantErr && err == nil {
				t.Fatal("ожидалась ошибка, получили nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("не ожидалась ошибка, получили: %v", err)
			}

			if tgMock.called != tt.wantSendCall {
				t.Fatalf("ожидали вызов SendNotification=%v, получили %v", tt.wantSendCall, tgMock.called)
			}

			if updaterMock.called != tt.wantUpdateCall {
				t.Fatalf("ожидали вызов UpdateNotificationStatus=%v, получили %v", tt.wantUpdateCall, updaterMock.called)
			}

			if tt.wantSendCall {
				if tgMock.chatID != validNotif.TgID {
					t.Fatalf("ожидали chatID=%d, получили %d", validNotif.TgID, tgMock.chatID)
				}
				if tgMock.text != validNotif.Text {
					t.Fatalf("ожидали text=%q, получили %q", validNotif.Text, tgMock.text)
				}
			}

			if tt.wantUpdateCall {
				if updaterMock.id != validNotif.ID {
					t.Fatalf("ожидали id=%d, получили %d", validNotif.ID, updaterMock.id)
				}
				if updaterMock.status != "sent" {
					t.Fatalf("ожидали status=%q, получили %q", "sent", updaterMock.status)
				}
			}
		})
	}
}
