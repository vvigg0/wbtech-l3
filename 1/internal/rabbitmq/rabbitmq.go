package myrabbitmq

import (
	"time"

	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type RabbitMQ struct {
	Client    *rabbitmq.RabbitClient
	Publisher *rabbitmq.Publisher
	Consumer  *rabbitmq.Consumer
}

func New(connURL string, handler rabbitmq.MessageHandler) (*RabbitMQ, error) {

	cfgRabbit := rabbitmq.ClientConfig{
		URL:            connURL,
		ConnectionName: "notifications",
		ConnectTimeout: 20 * time.Second,
		ReconnectStrat: retry.Strategy{
			Attempts: 5,
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
		return nil, err
	}
	if err := rabbitClient.DeclareExchange(
		"notifications.exchange",
		"direct",
		true,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	if err := rabbitClient.DeclareQueue(
		"notifications",
		"notifications.exchange",
		"notifications",
		true,
		false,
		false,
		nil); err != nil {
		return nil, err
	}

	consumer := rabbitmq.NewConsumer(
		rabbitClient,
		rabbitmq.ConsumerConfig{
			Queue:       "notifications",
			ConsumerTag: "notifications-consumer",
			Workers:     1},
		handler)

	publisher := rabbitmq.NewPublisher(rabbitClient, "notifications.exchange", "application/json")

	return &RabbitMQ{Client: rabbitClient, Publisher: publisher, Consumer: consumer}, nil
}
