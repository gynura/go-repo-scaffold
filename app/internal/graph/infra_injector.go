//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"strconv"
	"userservice/internal/infra/queue"
	"userservice/internal/infra/queue/rabbitmq"
)

type infraGraph struct {
	publisher queue.Publisher
}

var httpSet = wire.NewSet(
	wire.Struct(new(infraGraph), "*"),

	// RabbitMQ queue
	provideRabbitMQPublisherOptions,
	rabbitmq.ProvidePublisher,
	wire.Bind(new(queue.Publisher), new(*rabbitmq.Publisher)),
	//wire.Struct(new(rabbitmq.NewPublisherOptions), "*"),
)

func initializeInfraGraph(
	config Config,
) (*infraGraph, error) {
	wire.Build(httpSet)
	return &infraGraph{}, nil
}

func provideRabbitMQPublisherOptions(config Config) rabbitmq.NewPublisherOptions {
	port, err := strconv.Atoi(config.QueueConfig.RabbitMQ.Port)
	if err != nil {
		panic(err)
	}
	return rabbitmq.NewPublisherOptions{
		Username: config.QueueConfig.RabbitMQ.Username,
		Password: config.QueueConfig.RabbitMQ.Password,
		Host:     config.QueueConfig.RabbitMQ.Host,
		Port:     port,
		VHost:    config.QueueConfig.RabbitMQ.VHost,
	}
}
