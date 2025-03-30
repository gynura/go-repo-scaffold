// Package rabbitmq defines the infra layer for publishing messages through a rabbit mq queue
package rabbitmq

import (
	"encoding/json"
	"fmt"
	"net/url"
	"userservice/internal/core/events"
	"userservice/internal/infra/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

// userServiceExchangeName defines the only exchange in the user microservice domain
// A consumer must declare the queue and bind it to the exchange to receive messages
const userServiceExchangeName = "user_service"

type Publisher struct {
	channel *amqp.Channel
}

// NewPublisherOptions contain the attributes needed to construct a RabbitMQ connection
// This way who consumes this infra doesn't have to worry about constructing a proper RabbitMQ URL
type NewPublisherOptions struct {
	Username string
	Password string
	Host     string
	Port     int
	VHost    string
}

var _ queue.Publisher = new(Publisher)

// buildURL constructs the AMQP URL
func (opts NewPublisherOptions) buildURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		url.QueryEscape(opts.Username),
		url.QueryEscape(opts.Password),
		opts.Host,
		opts.Port,
		opts.VHost,
	)
}

// TODO: define custom rabbitmq infra errors so we can type errors at this layer
func ProvidePublisher(options NewPublisherOptions) (*Publisher, error) {
	conn, err := amqp.Dial(options.buildURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

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

func (r *Publisher) IsAlive() bool {
	// TODO: very dumbed down liveness
	if r.channel.IsClosed() {
		return false
	}

	return true
}
