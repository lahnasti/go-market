package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/go-market/internal/repository"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type Server struct {
	Db         repository.Repository
	ErrorChan  chan error
	deleteChan chan int
	Valid      *validator.Validate
	log        zerolog.Logger
	RabbitConn *amqp.Connection
	RabbitChan *amqp.Channel
}

func NewServer(ctx context.Context, db repository.Repository, zlog *zerolog.Logger) *Server {
	validate := validator.New()
	dChan := make(chan int, 5)
	errChan := make(chan error)
	srv := &Server{
		Db:         db,
		deleteChan: dChan,
		ErrorChan:  errChan,
		log:        *zlog,
		Valid:      validate,
	}
	err := srv.InitRabbit()
	if err != nil {
		zlog.Fatal().Err(err).Msg("RabbitMQ connection failed")
	}
	go srv.deleter(ctx)
	return srv
}
