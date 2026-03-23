package myrabbitmq

import (
	"log"
	"time"

	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type RabbitMQ struct {
	Publisher *rabbitmq.Publisher
	Consumer  *rabbitmq.Consumer
}

func New(connURL string, handler rabbitmq.MessageHandler) *RabbitMQ {

	cfgRabbit := rabbitmq.ClientConfig{
		URL:            connURL,
		ConnectionName: "notifications",
		ConnectTimeout: 10 * time.Second,
		ReconnectStrat: retry.Strategy{
			Attempts: 2,
			Delay:    1 * time.Second,
			Backoff:  2,
		},
		ProducingStrat: retry.Strategy{
			Attempts: 5,
			Delay:    1 * time.Second,
			Backoff:  2,
		},
		ConsumingStrat: retry.Strategy{
			Attempts: 5,
			Delay:    1 * time.Second,
			Backoff:  2,
		},
	}

	rabbitClient, err := rabbitmq.NewClient(cfgRabbit)
	if err != nil {
		log.Fatalf("ошибка создания клиента rabbitmq: %v", err)
	}

	if err := rabbitClient.DeclareExchange(
		"notifications.exchange",
		"direct",
		true,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("Ошибка создания exchange: %v", err)
	}

	if err := rabbitClient.DeclareQueue(
		"notifications",
		"notifications.exchange",
		"notifications",
		true,
		false,
		false,
		nil); err != nil {
		log.Fatalf("Ошибка создания очереди: %v", err)
	}

	consumer := rabbitmq.NewConsumer(
		rabbitClient,
		rabbitmq.ConsumerConfig{
			Queue:       "notifications",
			ConsumerTag: "notifications-consumer",
			Workers:     1},
		handler)

	publisher := rabbitmq.NewPublisher(rabbitClient, "notifications.exchange", "text/plain")

	return &RabbitMQ{Publisher: publisher, Consumer: consumer}
}
