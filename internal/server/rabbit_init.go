package server

import (
	"fmt"
	"os"

	amqp "github.com/streadway/amqp"
)

func (s *Server) InitRabbit() error {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_HOST")
	connStr := fmt.Sprintf("amqp://%s:%s@%s:5672/", "guest", "guest", rabbitURL)
	s.RabbitConn, err = amqp.Dial(connStr)
	if err != nil {
		return fmt.Errorf("failed to connect RabbitMQ: %w", err)
	}
	s.RabbitChan, err = s.RabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	return nil
}
