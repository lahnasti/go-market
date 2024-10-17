package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/auth/internal/server/responses"
	"github.com/lahnasti/go-market/common/models"
	"github.com/streadway/amqp"
)

func (s *Server) UserCheckHandler(msg amqp.Delivery, ctx *gin.Context) {
	var userCheckMsg models.UserCheckMessage
	if err := json.Unmarshal(msg.Body, &userCheckMsg); err != nil {
		s.log.Error().Err(err).Msg("Failed to unmarshal user check message")
		responses.SendError(ctx, http.StatusBadRequest, "Failed to unmarshal user check message", err)
		return
	}
	user, err := s.Db.GetUserProfile(userCheckMsg.UserID)
	response := models.UserCheckResponse{
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

func (s *Server) SendUserCheckResponse(response models.UserCheckResponse) error {
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
