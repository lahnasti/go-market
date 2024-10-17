package server

import (
	"encoding/json"
	"fmt"

	"github.com/lahnasti/go-market/common/models"
)

func (s *Server) WaitForUserCheckResponse() (models.UserCheckResponse, error) {
	var response models.UserCheckResponse

	// Потребляем сообщения из очереди ответов
	messages, err := s.Rabbit.Channel.Consume(
		"user_check_response_queue", // Имя очереди с ответами
		"",                          // consumer tag
		true,                        // auto-ack, автоматическое подтверждение
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // аргументы
	)
	if err != nil {
		return response, fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Ожидаем получения сообщения из очереди
	for msg := range messages {
		if err := json.Unmarshal(msg.Body, &response); err != nil {
			return response, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// Проверяем валидность ответа
		if response.Valid {
			return response, nil
		}

	}

	// Возвращаем ошибку, если не получен ответ
	return response, fmt.Errorf("no response received")
}
