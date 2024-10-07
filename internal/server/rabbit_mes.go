package server

import (
	"encoding/json"
	"fmt"

	"github.com/lahnasti/go-market/internal/models"
	amqp "github.com/streadway/amqp"
)

type UserRegisteredMes struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Server) sendUserRegisteredMessage(user models.User, userID int) error {
	// Создаем сообщение
	message := UserRegisteredMes{
		UserID:   userID,
		Username: user.Username,
		Email:    user.Email,
	}

	// Преобразуем сообщение в JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Публикуем сообщение в очередь
	err = s.RabbitChan.Publish(
		"",           // exchange
		"user_queue", // routing key (название очереди)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
