package rabbitmq

import (
	"fmt"
	"os"

	amqp "github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
}
func InitRabbit() (*RabbitMQ, error) {

	rabbitURL := os.Getenv("RABBITMQ_HOST")
	if rabbitURL == "" {
		rabbitURL = "localhost"
	}

	rabbitPort := os.Getenv("RABBITMQ_PORT")
	if rabbitPort == "" {
		rabbitPort = "5672"
	}

	connStr := fmt.Sprintf("ampq://%s:%s@%s:%s/", "guest", "guest", rabbitURL, rabbitPort)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect RabbitMQ: %w", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil,  fmt.Errorf("failed to open a channel: %w", err)
	}
	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}, nil
}

func (r *RabbitMQ) CloseRabbit() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}
