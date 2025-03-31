// Package rabbitmq defines the infra layer for publishing messages through a rabbit mq queue
package rabbitmq

import (
	"context"
	"fmt"
	"net/url"
	"userservice/internal/core/events"
	"userservice/internal/infra/queue"
	"userservice/internal/pkg/commons/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

// userServiceExchangeName defines the only exchange in the user microservice domain
// A consumer must declare the queue and bind it to the exchange to receive messages
// TODO: make this configurable
const userServiceExchangeName = "user_service"

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	closeCh chan *amqp.Error
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
		conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		userServiceExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	publisher := &Publisher{
		conn:    conn,
		channel: ch,
		closeCh: make(chan *amqp.Error),
	}

	// We need to listen for uncontrolled channel closure notifications
	ch.NotifyClose(publisher.closeCh)
	go func() {
		err := <-publisher.closeCh
		if err != nil {
			logger.NewLogger().Error(context.Background()).Message(err.Error())
		}
	}()

	return publisher, nil
}

func (r *Publisher) Publish(event events.ServiceEvent) error {
	body, err := event.Payload()
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

// TODO: gracefully close this when the whole server closes
func (r *Publisher) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
