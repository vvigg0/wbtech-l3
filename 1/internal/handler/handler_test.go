package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/dto"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/repository"
)

type serviceMock struct {
	createFn func(ctx context.Context, req dto.CreateNotificationsRequest) ([]int, error)
	getFn    func(id int) (*dto.NotificationStatusDTO, error)
	cancelFn func(id int) error
}

func (m *serviceMock) CreateNotifications(ctx context.Context, req dto.CreateNotificationsRequest) ([]int, error) {
	return m.createFn(ctx, req)
}

func (m *serviceMock) GetNotificationStatus(id int) (*dto.NotificationStatusDTO, error) {
	return m.getFn(id)
}

func (m *serviceMock) CancelNotification(id int) error {
	return m.cancelFn(id)
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/notify", h.CreateNotification)
	r.GET("/notify/:id", h.GetNotificationStatus)
	r.DELETE("/notify/:id", h.CancelNotification)
	return r
}

func TestCreateNotification_OK(t *testing.T) {
	svc := &serviceMock{
		createFn: func(ctx context.Context, req dto.CreateNotificationsRequest) ([]int, error) {
			return []int{1, 2}, nil
		},
		getFn: func(id int) (*dto.NotificationStatusDTO, error) {
			return nil, nil
		},
		cancelFn: func(id int) error {
			return nil
		},
	}

	h := New(svc)
	r := setupRouter(h)

	body := []byte(`{
		"notifs": [
			{
				"text": "hello",
				"telegram_ID": 123,
				"send_at": "2026-03-26T12:00:00Z"
			}
		]
	}`)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateNotification_BadJSON(t *testing.T) {
	svc := &serviceMock{
		createFn: func(ctx context.Context, req dto.CreateNotificationsRequest) ([]int, error) {
			return nil, nil
		},
		getFn: func(id int) (*dto.NotificationStatusDTO, error) {
			return nil, nil
		},
		cancelFn: func(id int) error {
			return nil
		},
	}

	h := New(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("ожидали 400, получили %d", w.Code)
	}
}

func TestGetNotificationStatus_OK(t *testing.T) {
	svc := &serviceMock{
		createFn: func(ctx context.Context, req dto.CreateNotificationsRequest) ([]int, error) {
			return nil, nil
		},
		getFn: func(id int) (*dto.NotificationStatusDTO, error) {
			return &dto.NotificationStatusDTO{
				ID:     id,
				Status: "queued",
			}, nil
		},
		cancelFn: func(id int) error {
			return nil
		},
	}

	h := New(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/notify/10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json decode error: %v", err)
	}
}

func TestDeleteNotification_NotFound(t *testing.T) {
	svc := &serviceMock{
		createFn: func(ctx context.Context, req dto.CreateNotificationsRequest) ([]int, error) {
			return nil, nil
		},
		getFn: func(id int) (*dto.NotificationStatusDTO, error) {
			return nil, nil
		},
		cancelFn: func(id int) error {
			return repository.ErrNoNotification
		},
	}

	h := New(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/notify/999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("ожидали 404, получили %d, body=%s", w.Code, w.Body.String())
	}
}
