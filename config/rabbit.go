package config

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateRabbitMQClient(config *RabbitMQConfig) (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.Username, config.Password, config.Host, config.Port)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}
