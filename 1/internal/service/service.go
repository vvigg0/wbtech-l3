package service

import (
	myrabbitmq "github.com/vvigg0/wbtech-l3/l3/1/internal/rabbitmq"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/repository"
)

type Service struct {
	repo   *repository.Repository
	Rabbit *myrabbitmq.RabbitMQ
}

func New(r *repository.Repository, queue *myrabbitmq.RabbitMQ) *Service {
	return &Service{r, queue}
}
