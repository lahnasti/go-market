package server

import (
	"encoding/json"
	"fmt"

	"github.com/lahnasti/go-market/auth/internal/models"
)

func (s *Server) sendUserRegisteredMessage(user models.User, id int) error {
	// Подготавливаем сообщение
	message := map[string]interface{}{
		"userID":   id,
		"username": user.Username,
		"email":    user.Email,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal user registration message: %w", err)
	}

	// Публикуем сообщение в очередь RabbitMQ
	err = s.Rabbit.PublishMessage("user_registered_queue", messageBytes)
	if err != nil {
		return fmt.Errorf("failed to publish user registration message: %w", err)
	}

	return nil
}
