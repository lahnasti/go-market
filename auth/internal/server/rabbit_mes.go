package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/auth/internal/server/responses"
	"github.com/streadway/amqp"
)

func (s *Server) UserCheckHandler(msg amqp.Delivery, ctx *gin.Context) {
	var userCheckMsg UserCheckMessage
	if err := json.Unmarshal(msg.Body, &userCheckMsg); err != nil {
		s.log.Error().Err(err).Msg("Failed to unmarshal user check message")
		responses.SendError(ctx, http.StatusBadRequest, "Failed to unmarshal user check message", err)
		return
	}
	user, err := s.Db.GetUserProfile(userCheckMsg.UserID)
	response := UserCheckResponse{
		Valid: false,
	}
	if err == nil {
		response.Valid = true
		response.User = user
	} else {
		s.log.Error().Err(err).Msg("Failed to get user profile")
	}
	err = s.SendUserCheckResponse(response)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to send user check response")
		responses.SendError(ctx, http.StatusInternalServerError, "Failed to send user check response", err)
	}
}

func (s *Server) SendUserCheckResponse(response UserCheckResponse) error {
	body, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	err = s.Rabbit.Channel.Publish(
		"",
		"user_check_response_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish response: %w", err)
	}
	return nil
}

func (s *Server) WaitForUserCheckResponse() (UserCheckResponse, error) {
	var response UserCheckResponse

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
