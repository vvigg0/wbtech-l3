package handler

import (
	"github.com/vvigg0/wbtech-l3/l3/1/internal/service"
)

type Handler struct {
	srvc *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{srvc: s}
}
