package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/go-market/auth/internal/repository"
	"github.com/lahnasti/go-market/lib/rabbitmq"
	"github.com/rs/zerolog"
)

type Server struct {
	Db         repository.UserRepository
	ErrorChan  chan error
	deleteChan chan int
	Valid      *validator.Validate
	log        zerolog.Logger
	Rabbit     *rabbitmq.RabbitMQ
}

func NewServer(ctx context.Context, db repository.UserRepository, zlog *zerolog.Logger) *Server {
	validate := validator.New()
	errChan := make(chan error)
	rabbitClient, err := rabbitmq.InitRabbit()
	if err != nil {
		zlog.Fatal().Err(err).Msg("RabbitMQ connection failed")
	}
	srv := &Server{
		Db:        db,
		ErrorChan: errChan,
		log:       *zlog,
		Valid:     validate,
		Rabbit: rabbitClient,
	}
	return srv
}

func (s *Server) Close() {
	if s.Rabbit != nil {
		s.Rabbit.CloseRabbit()
	}
}