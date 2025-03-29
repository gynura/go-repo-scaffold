// Package rabbitmq defines the infra layer forp publishing messages through a rabbit mq queue
package rabbitmq

import (
	"encoding/json"
	"userservice/internal/core/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

// userServiceExchangeName defines the only exchange in the user microservice domain
const userServiceExchangeName = "user_service"

type Publisher struct {
	channel *amqp.Channel
}

// TODO: define custom rabbitmq infra errors so we can type errors at this layer
func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{channel: ch}, nil
}

func (r *Publisher) Publish(event events.ServiceEvent) error {
	body, err := json.Marshal(event.Payload())
	if err != nil {
		return err
	}

	return r.channel.Publish(
		userServiceExchangeName,
		event.EventType().String(),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
